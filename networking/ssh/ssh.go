package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func DialSSHClient(host string, port string, user string, password string) (*ssh.Client, error) {

	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
