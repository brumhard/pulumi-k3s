package k3s

import (
	"bytes"
	"fmt"
	"io"
	"path"

	"github.com/brumhard/pulumi-k3s/provider/pkg/sshexec"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RemoteExecutor struct {
	fileHandler *sftp.Client
	cmdHandler  *sshexec.Client
	addr        string
	sudoPrefix  string
}

func NewExecutorForNode(node Node, useSudo bool) (*RemoteExecutor, error) {
	return NewRemoteExecutor(fmt.Sprintf("%s:22", node.Host), node.User, []byte(node.PrivateKey), useSudo)
}

func NewRemoteExecutor(addr, user string, sshkeyPEM []byte, useSudo bool) (*RemoteExecutor, error) {
	signer, err := ssh.ParsePrivateKey(sshkeyPEM)
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	fileHandler, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	cmdHandler := sshexec.NewClientFromSSH(conn)

	sudoPrefix := ""
	if useSudo {
		sudoPrefix = "sudo "
	}

	return &RemoteExecutor{
		fileHandler: fileHandler,
		cmdHandler:  cmdHandler,
		addr:        addr,
		sudoPrefix:  sudoPrefix,
	}, nil
}

func (e *RemoteExecutor) ExecuteScript(script []byte) error {
	tmpScriptLocation := "/tmp/pulumi-k3s-tmp-script"

	file, err := e.fileHandler.Create(tmpScriptLocation)
	if err != nil {
		return errors.Wrapf(err, "failed to create tmp file for script")
	}

	defer func() {
		file.Close()
		e.fileHandler.Remove(tmpScriptLocation)
	}()

	if _, err := io.Copy(file, bytes.NewReader(script)); err != nil {
		return errors.Wrap(err, "failed to write to tmp script file")
	}

	executeScript := &sshexec.Cmd{Command: e.sudoSprintF(
		"sh %s", tmpScriptLocation,
	)}

	if err := e.cmdHandler.Run(executeScript); err != nil {
		return errors.Wrap(err, "failed to execute script")
	}

	return nil
}

func (e *RemoteExecutor) CopyFile(fileReader io.Reader, dest string) error {
	tmpFilePath := path.Join("/tmp", path.Base(dest))

	file, err := e.fileHandler.Create(tmpFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to create tmp file at %s", tmpFilePath)
	}

	defer file.Close()

	if _, err := io.Copy(file, fileReader); err != nil {
		return errors.Wrapf(err, "failed to write to %s", tmpFilePath)
	}

	mvCmd := &sshexec.Cmd{Command: e.sudoSprintF(
		"mv %s %s", tmpFilePath, dest,
	)}
	if err := e.cmdHandler.Run(mvCmd); err != nil {
		return errors.Wrapf(err, "failed to move from %s to %s", tmpFilePath, dest)
	}

	return nil
}

func (e *RemoteExecutor) SudoCombinedOutput(cmd string) (string, error) {
	return e.CombinedOutput(e.sudoSprintF(cmd))
}

func (e *RemoteExecutor) CombinedOutput(cmd string) (string, error) {
	stdouterr := &bytes.Buffer{}
	err := e.cmdHandler.Run(&sshexec.Cmd{
		Command: cmd,
		Stderr:  stdouterr,
		Stdout:  stdouterr,
	})
	if err != nil {
		return "", errors.Wrap(errors.Wrap(err, stdouterr.String()), e.addr)
	}

	return stdouterr.String(), nil
}

func (e *RemoteExecutor) SudoOutput(cmd string) (string, error) {
	return e.Output(e.sudoSprintF(cmd))
}

func (e *RemoteExecutor) Output(cmd string) (string, error) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	err := e.cmdHandler.Run(&sshexec.Cmd{
		Command: cmd,
		Stderr:  stderr,
		Stdout:  stdout,
	})
	if err != nil {
		return "", errors.Wrap(errors.Wrap(err, stderr.String()), e.addr)
	}

	return stdout.String(), nil
}

func (e *RemoteExecutor) sudoSprintF(format string, a ...interface{}) string {
	format = e.sudoPrefix + format
	return fmt.Sprintf(format, a...)
}
