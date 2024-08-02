package remote

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ontime-dev/ont/internal/escape"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func Run(user, server, cmd, prvtKeyFile, pubKeyFile string, stat, verbose bool) error {

	if !stat {
		cmd = fmt.Sprintf("bash %s", cmd)
	}
	userSSHDir := fmt.Sprintf("/home/%s/.ssh", user)

	if verbose {
		escape.LogPrintf("DEBUG(RMTRN): Getting ssh private key for user %s", user)
	}
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

	if verbose {
		escape.LogPrintf("DEBUG(RMTRN): Getting IP and hostname of the server %s", server)
	}
	serverIP, serverHost, err := getServerIPHOST(server, verbose)
	if err != nil {
		return err
	}

	if verbose {
		escape.LogPrintf("DEBUG(RMTRN): Trying with server IP '%s'", serverIP[0])
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(serverIP[0].String(), "22"), config)
	if err != nil {
		if verbose {
			escape.LogPrintf("DEBUG(RMTRN): Trying with server hostname '%s'", serverHost[0])
		}
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
	if verbose {
		escape.LogPrint("DEBUG(RMTRN): Connetion Succeeded. Proceeding")
	}

	session, err := client.NewSession()

	if err != nil {
		return err
	}

	err = session.Run(cmd)

	if verbose {
		escape.LogPrint("DEBUG(RMTRN): Closing open ssh sessions")
	}
	session.Close()
	client.Close()

	return err
}

func getServerIPHOST(server string, verbose bool) ([]net.IP, []string, error) {

	var serverHost []string
	var serverIP []net.IP
	var err error

	ip := net.ParseIP(server)
	//If ip == nil, it is a hostname, if not nil, it is an ip.
	if ip != nil {
		if verbose {
			escape.LogPrintf("DEBUG(GTSRVRIPHST): Reverse lookup the IP %s", server)
		}
		serverHost, err = net.LookupAddr(server)
		serverIP = []net.IP{ip}
		if err != nil {
			return nil, nil, err
		}
	} else {
		if verbose {
			escape.LogPrintf("DEBUG(GTSRVRIPHST): Resolving the server hostname %s", server)
		}
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
