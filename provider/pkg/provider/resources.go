package provider

import (
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func makeOrUpdateCluster(name string, configMap map[string]interface{}) (*Cluster, error) {
	configBytes, err := yaml.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	cluster := Cluster{
		User: "root", //default
	}
	if err := yaml.Unmarshal(configBytes, &cluster); err != nil {
		return nil, err
	}

	if err := cluster.Validate(); err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp("", name)
	if err != nil {
		return nil, err
	}

	sshKeyPath, kubeconfigPath := path.Join(tempDir, "sshkey"), path.Join(tempDir, "kubeconfig")
	if err := os.WriteFile(sshKeyPath, []byte(cluster.PrivateKey), os.ModePerm); err != nil {
		return nil, err
	}

	_, err = exec.Command(
		"k3sup", "install",
		"--ip", cluster.IP,
		"--user", cluster.User,
		"--ssh-key", sshKeyPath,
		"--local-path", kubeconfigPath,
	).CombinedOutput()
	if err != nil {
		return nil, err
	}

	kubeconfigBytes, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	cluster.KubeConfig = string(kubeconfigBytes)
	return &cluster, nil
}

type Cluster struct {
	IP         string `json:"ip" yaml:"ip"`
	User       string `json:"user" yaml:"user"`
	KubeConfig string `json:"kubeconfig" yaml:"kubeconfig"`
	PrivateKey string `json:"privateKey" yaml:"privateKey"`
}

var (
	ErrRequiredProperty = errors.New("property is required")
	ErrOutputOnly       = errors.New("property is only for output")
)

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
