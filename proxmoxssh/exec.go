package proxmoxssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
)

func SimpleRemoteRun(cfg Config, cmd string, stdininput string) (string, error) {
	var authmethods []ssh.AuthMethod

	if cfg.PrivateKey != nil && len(cfg.PrivateKey) > 0 {
		// Prepare key
		key, err := ssh.ParsePrivateKey(cfg.PrivateKey)
		if err != nil {
			return "", errors.Wrap(err, "parse private key")
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
			return "", errors.Wrap(err, "parse public key")
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
		return "", errors.Wrap(err, "dial")
	}

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "create new session")
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
		return "", errors.Wrap(err, "run command")
	}

	// Cleanup
	err = session.Close()
	if err != nil && err != io.EOF {
		return "", errors.Wrap(err, "closing session")
	}
	err = client.Close()
	if err != nil {
		return "", errors.Wrap(err, "closing client")
	}

	// End
	return b.String(), nil
}
