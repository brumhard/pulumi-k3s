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

// TODO: generate schema from this struct
type Cluster struct {
	MasterNodes []Node `json:"masterNodes"`
	Agents      []Node `json:"agents"`
	KubeConfig  string `json:"kubeconfig"`
}

// TODO: user and privatekey in provider as defaults (like region in openstack)
type Node struct {
	Host       string `json:"host"`
	User       string `json:"user"`
	PrivateKey string `json:"privateKey"`
}

func MakeOrUpdateCluster(name string, cluster *Cluster) error {
	if err := cluster.Validate(); err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", name)
	if err != nil {
		return err
	}

	kubeconfig, err := setupMasterNode(cluster.MasterNodes[0], tempDir)
	if err != nil {
		return err
	}

	cluster.KubeConfig = kubeconfig

	return nil
}

func setupMasterNode(node Node, tempDir string) (string, error) {
	sshKeyPath, kubeconfigPath := path.Join(tempDir, "sshkey"), path.Join(tempDir, "kubeconfig")
	if err := os.WriteFile(sshKeyPath, []byte(node.PrivateKey), os.ModePerm); err != nil {
		return "", err
	}

	_, err := exec.Command(
		"k3sup", "install",
		// TODO: should --host be used here
		"--ip", node.Host,
		"--user", node.User,
		"--ssh-key", sshKeyPath,
		"--local-path", kubeconfigPath,
	).CombinedOutput()
	if err != nil {
		return "", err
	}

	kubeconfigBytes, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return "", err
	}

	return string(kubeconfigBytes), nil
}

func (c Cluster) Validate() error {
	if len(c.MasterNodes) != 1 {
		return errors.New("only clusters with exactly 1 master node supported")
	}

	if len(c.Agents) > 0 {
		return errors.New("agents are currently not supported")
	}

	if c.KubeConfig != "" {
		return errors.Wrap(ErrOutputOnly, "kubeconfig")
	}

	for _, n := range append(c.MasterNodes, c.Agents...) {
		if err := n.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func DeleteCluster(cluster *Cluster) error {
	for _, n := range cluster.MasterNodes {
		// TODO: handle error if already gone
		if err := executeOnNode(n, "/usr/local/bin/k3s-uninstall.sh"); err != nil {
			return err
		}
	}

	for _, n := range cluster.Agents {
		if err := executeOnNode(n, "/usr/local/bin/k3s-agent-uninstall.sh"); err != nil {
			return err
		}
	}

	return nil
}

func executeOnNode(node Node, command string) error {
	client, err := sshexec.NewClient(fmt.Sprintf("%s:22", node.Host), node.User, []byte(node.PrivateKey))
	if err != nil {
		return err
	}

	stderr := &bytes.Buffer{}
	err = client.Run(&sshexec.Cmd{
		Command: command,
		Stderr:  stderr,
	})
	if err != nil {
		return errors.Wrap(errors.Wrap(err, stderr.String()), node.Host)
	}

	return nil
}

func (n Node) Validate() error {
	if n.Host == "" {
		return errors.Wrap(ErrRequiredProperty, "host")
	}

	if n.PrivateKey == "" {
		return errors.Wrap(ErrRequiredProperty, "privateKey")
	}

	if n.User == "" {
		return errors.Wrap(ErrRequiredProperty, "user")
	}

	return nil
}
