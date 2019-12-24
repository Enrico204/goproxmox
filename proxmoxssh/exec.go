package proxmoxssh

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

func SimpleRemoteRun(cfg Config, cmd string, stdininput string) (string, error) {
	var authmethods []ssh.AuthMethod

	if cfg.PrivateKey != nil && len(cfg.PrivateKey) > 0 {
		// Prepare key
		key, err := ssh.ParsePrivateKey(cfg.PrivateKey)
		if err != nil {
			return "", err
		}
		authmethods = append(authmethods, ssh.PublicKeys(key))
	}

	if cfg.Password != "" {
		authmethods = append(authmethods, ssh.Password(cfg.Password))
	}

	var hostKeyCallback = ssh.InsecureIgnoreHostKey()
	if cfg.HostPublicKey != nil {
		hostkey, err := ssh.ParsePublicKey(cfg.HostPublicKey)
		if err != nil {
			return "", err
		}
		hostKeyCallback = ssh.FixedHostKey(hostkey)
	}

	// Authentication
	config := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authmethods,
		HostKeyCallback: hostKeyCallback,
	}

	// Connect
	client, err := ssh.Dial("tcp", net.JoinHostPort(cfg.Hostname, fmt.Sprint(cfg.Port)), config)
	if err != nil {
		return "", err
	}

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	// Retrieve the output
	var b bytes.Buffer
	session.Stdout = &b
	if stdininput != "" {
		session.Stdin = bytes.NewBufferString(stdininput)
	}

	// Run the command
	err = session.Run(cmd)
	if err != nil {
		return "", err
	}

	// Cleanup
	err = session.Close()
	if err != nil {
		return "", err
	}
	err = client.Close()
	if err != nil {
		return "", err
	}

	// End
	return b.String(), nil
}
