package remote

import (
	"errors"
	"fmt"
	"net"
	"ont/internal/escape"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func addHostKey(host, user string, remote net.Addr, pubKey ssh.PublicKey) error {
	// add host key if host is not found in known_hosts, error object is return, if nil then connection proceeds,
	// if not nil then connection stops.
	//khFilePath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	knownHostsFilepath := fmt.Sprintf("/home/%s/.ssh/known_hosts", user)

	file, err := os.OpenFile(knownHostsFilepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		escape.LogPrint(err.Error())
	}
	defer file.Close()

	knownHosts := knownhosts.Normalize(remote.String())
	_, err = file.WriteString(knownhosts.Line([]string{knownHosts}, pubKey))
	return err
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

func Run(user, server, cmd, prvt_key string) error {
	cmd = fmt.Sprintf("bash %s", cmd)

	user_ssh_dir := fmt.Sprintf("/home/%s/.ssh", user)

	// allowed_hosts_file := fmt.Sprintf("%s/known_hosts", user_ssh_dir)
	// publicKeyBytes, _ := os.ReadFile(allowed_hosts_file)
	// _, hosts, pubkey, _, _, _ := ssh.ParseKnownHosts(publicKeyBytes)
	// escape.LogPrint(pubkey)
	// escape.LogPrint(hosts)

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

	var keyErr *knownhosts.KeyError

	config := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
			kh := checkKnownHosts(user)

			err := kh(host, remote, pubKey)

			if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {

				return keyErr
			} else if errors.As(err, &keyErr) && len(keyErr.Want) == 0 {

				return addHostKey(host, user, remote, pubKey)
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
