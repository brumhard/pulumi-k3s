package sshexec

import (
	"io"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	conn *ssh.Client
}

func NewClient(address, user string, sshkeyPEM []byte) (*Client, error) {
	signer, err := ssh.ParsePrivateKey(sshkeyPEM)
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", address, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

type Cmd struct {
	Command string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
}

func (c *Client) Run(cmd *Cmd) error {
	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdin = cmd.Stdin
	session.Stderr = cmd.Stderr
	session.Stdout = cmd.Stdout

	return session.Run(cmd.Command)
}
