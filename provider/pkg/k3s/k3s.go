package k3s

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/brumhard/pulumi-k3s/provider/pkg/sshexec"
	"github.com/pkg/errors"
)

var (
	ErrRequiredProperty = errors.New("property is required")
	ErrOutputOnly       = errors.New("property is only for output")
)

type Cluster struct {
	IP         string `json:"ip" yaml:"ip"`
	User       string `json:"user" yaml:"user"`
	KubeConfig string `json:"kubeconfig" yaml:"kubeconfig"`
	PrivateKey string `json:"privateKey" yaml:"privateKey"`
}

func MakeOrUpdateCluster(name string, cluster *Cluster) error {
	if err := cluster.Validate(); err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", name)
	if err != nil {
		return err
	}

	sshKeyPath, kubeconfigPath := path.Join(tempDir, "sshkey"), path.Join(tempDir, "kubeconfig")
	if err := os.WriteFile(sshKeyPath, []byte(cluster.PrivateKey), os.ModePerm); err != nil {
		return err
	}

	_, err = exec.Command(
		"k3sup", "install",
		"--ip", cluster.IP,
		"--user", cluster.User,
		"--ssh-key", sshKeyPath,
		"--local-path", kubeconfigPath,
	).CombinedOutput()
	if err != nil {
		return err
	}

	kubeconfigBytes, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return err
	}

	cluster.KubeConfig = string(kubeconfigBytes)
	return nil
}

func DeleteCluster(cluster *Cluster) error {
	client, err := sshexec.NewClient(fmt.Sprintf("%s:22", cluster.IP), cluster.User, []byte(cluster.PrivateKey))
	if err != nil {
		return err
	}

	stderr := &bytes.Buffer{}
	err = client.Run(&sshexec.Cmd{
		Command: "/usr/local/bin/k3s-uninstall.sh",
		Stderr:  stderr,
	})
	if err != nil {
		return errors.Wrap(err, stderr.String())
	}

	return nil
}

func (c Cluster) Validate() error {
	if c.IP == "" {
		return errors.Wrap(ErrRequiredProperty, "ip")
	}

	if c.PrivateKey == "" {
		return errors.Wrap(ErrRequiredProperty, "privateKey")
	}

	if c.KubeConfig != "" {
		return errors.Wrap(ErrOutputOnly, "kubeconfig")
	}

	return nil
}
