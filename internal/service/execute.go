package service

import (
	"database/sql"
	"fmt"
	"net"
	"ont/internal/dbopts"
	"ont/internal/escape"
	"ont/internal/run"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func Execute(db *sql.DB, table string, job dbopts.Jobs) {
	//Enable below when verbose
	//escape.LogPrint("Executing script", job.Script)
	escape.LogPrintf("Executing job with ID %d for user '%s' on %s", job.Id, table, job.RunOn)
	uid, homeDir := getUserInfo(table)
	cmd := exec.Command(job.Script)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid}
	cmd.Dir = homeDir
	go func() {
		if job.RunOn == "local" {
			escape.LogPrint(job.RunOn)
			err := cmd.Run()
			if err != nil {

				escape.LogPrint(err.Error())
			}
		} else {
			escape.LogPrint(job.RunOn)
			err := remoteRun(table, job.RunOn, job.Script, "")
			if err != nil {
				escape.LogPrint(err.Error())
			}
		}
	}()
	err := ChangeExecTime(db, table, job)
	if err != nil {
		escape.LogPrint(err.Error())
	}
}

func ChangeExecTime(db *sql.DB, table string, job dbopts.Jobs) error {
	crntTime, err := time.Parse("15:04:05 Jan 02 2006", job.Exec_time)
	if err != nil {
		return err
	}

	number, last_char := run.GetLastChar(job.Every)
	job.Exec_time, err = run.ParseEvery(crntTime, last_char, number)
	if err != nil {
		return err
	}

	//Enable below when verbose
	//escape.LogPrint("Changing time to: ", job.Exec_time)
	_, err = dbopts.Insert(db, table, job, false)

	if err != nil {
		return err
	}
	return nil

}

func getUserInfo(table string) (uint32, string) {
	user, err := user.Lookup(table)
	if err != nil {
		escape.LogPrint(err)
	}
	userUID, err := strconv.ParseUint(user.Uid, 10, 32)
	if err != nil {
		escape.LogPrint(err)
	}
	return uint32(userUID), user.HomeDir
}

func createKnownHosts(user string) string {
	knownHostsFilepath := fmt.Sprintf("/home/%s/.ssh/known_hosts", user)
	file, err := os.OpenFile(knownHostsFilepath, os.O_CREATE, 0600)
	if err != nil {
		escape.LogPrint(err.Error())
	}
	file.Close()

	return knownHostsFilepath
}

func checkKnownHosts(user string) ssh.HostKeyCallback {
	knownHostsFile := createKnownHosts(user)
	knownHosts, err := knownhosts.New(knownHostsFile)

	if err != nil {
		escape.LogPrint(err.Error())
		return nil
	}
	return knownHosts
}

func remoteRun(user, server, cmd, prvt_key string) error {
	cmd = fmt.Sprintf("bash %s", cmd)

	user_ssh_dir := fmt.Sprintf("/home/%s/.ssh", user)

	allowed_hosts_file := fmt.Sprintf("%s/known_hosts", user_ssh_dir)
	publicKeyBytes, _ := os.ReadFile(allowed_hosts_file)
	_, hosts, pubkey, _, _, _ := ssh.ParseKnownHosts(publicKeyBytes)
	escape.LogPrint(pubkey)
	escape.LogPrint(hosts)

	if prvt_key == "" {
		prvt_key = fmt.Sprintf("%s/id_rsa", user_ssh_dir)
	}
	privateKey, err := os.ReadFile(prvt_key)

	if err != nil {
		return err
	}
	key, err := ssh.ParsePrivateKey(privateKey)

	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
			kh := checkKnownHosts(user)
			err := kh(host, remote, pubKey)
			if err != nil {
				return err
			}
			return nil
		}),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(server, "22"), config)

	if err != nil {
		return err
	}

	session, err := client.NewSession()

	if err != nil {
		return err
	}

	defer session.Close()

	err = session.Run(cmd)
	if err != nil {
		return err
	}
	return err

}
