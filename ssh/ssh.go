package ssh

import (
	"fmt"
	"time"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

var (
	user string
)

func User(u string) {
	user = u
}

func init() {
	User("root")
}

func Run(key []byte, host, cmd string) (string, error) {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("error parsing key %w", err)
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 15,
	}
	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return "", fmt.Errorf("error connecting to host %s %w", host, err)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("error creating session %w", err)
	}
	defer session.Close()
	b, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("error running commmand %s %w", cmd, err)
	}
	return string(b), nil
}

func CopyTo(key []byte, host, source, dest string) error {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("error parsing key %w", err)
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 15,
	}
	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return fmt.Errorf("error connecting to host %s %w", host, err)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error creating session %w", err)
	}
	defer session.Close()
	if err := scp.CopyPath(source, dest, session); err != nil {
		return fmt.Errorf("error copying file to %s", host)
	}
	return nil
}
