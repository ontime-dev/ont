package remote

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func Run(user, server, cmd, prvtKeyFile, pubKeyFile string, stat bool) error {

	if !stat {
		cmd = fmt.Sprintf("bash %s", cmd)
	}
	userSSHDir := fmt.Sprintf("/home/%s/.ssh", user)

	prvtKey, err := getPrivateKey(userSSHDir, prvtKeyFile)
	if err != nil {
		return err
	}

	hostcallback, err := newHostKeyCallBackFunc(userSSHDir)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostcallback,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(prvtKey),
		},
	}

	serverIP, serverHost, err := getServerIPHOST(server)
	if err != nil {
		return err
	}

	//verbose
	//escape.LogPrint("Trying with server IP")
	client, err := ssh.Dial("tcp", net.JoinHostPort(serverIP[0].String(), "22"), config)
	if err != nil {
		//verbose
		// escape.LogPrint("Trying with server hostname")
		client, err = ssh.Dial("tcp", net.JoinHostPort(serverHost[0], "22"), config)
		if err != nil {
			if strings.Contains(err.Error(), "knownhosts: key is unknown") {
				return errors.New("error: ssh handshake failed. Make sure the remote server is already to the Known_hosts file on the ontd server")
			} else if strings.Contains(err.Error(), "ssh: unable to authenticate") {
				return errors.New("error: ssh handshake failed. Make sure that the PubKey authentication is enabled for the user")
			} else {
				return err
			}
		}
	}
	//verbose
	// escape.LogPrint("Connection Succeeded. Proceeding.")

	session, err := client.NewSession()

	if err != nil {
		return err
	}

	err = session.Run(cmd)

	session.Close()
	client.Close()

	return err
}

func getServerIPHOST(server string) ([]net.IP, []string, error) {

	var serverHost []string
	var serverIP []net.IP
	var err error

	ip := net.ParseIP(server)
	//If ip == nil, it is a hostname, if not nil, it is an ip.
	if ip != nil {
		//verbose
		//escape.LogPrint(Reverse resolving the IP)
		serverHost, err = net.LookupAddr(server)
		serverIP = []net.IP{ip}
		if err != nil {
			return nil, nil, err
		}
	} else {
		//verbose
		//escape.LogPrint(Resolving the server")
		serverIP, err = net.LookupIP(server)
		serverHost = []string{server}
		if err != nil {
			return nil, nil, err
		}
	}
	return serverIP, serverHost, err
}

func getPrivateKey(userSShDir, prvtKeyFile string) (ssh.Signer, error) {

	if prvtKeyFile == "" {
		prvtKeyFile = fmt.Sprintf("%s/id_rsa", userSShDir)
	}

	privateKey, err := os.ReadFile(prvtKeyFile)
	if err != nil {
		return nil, fmt.Errorf("%s. RSA only is allowed at the moment", err.Error())
	}

	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return key, err
}

func newHostKeyCallBackFunc(userSShDir string) (ssh.HostKeyCallback, error) {

	knownhostsFile := fmt.Sprintf("%s/known_hosts", userSShDir)

	hostcallback, err := knownhosts.New(knownhostsFile)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}

	return hostcallback, err
}
