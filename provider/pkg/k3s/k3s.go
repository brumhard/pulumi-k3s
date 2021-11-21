package k3s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/brumhard/pulumi-k3s/provider/pkg/sshexec"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
)

const (
	getScript              = "curl -sfL https://get.k3s.io"
	useSudo                = true
	channelURL             = "https://update.k3s.io/v1-release/channels"
	autoDeployManifestPath = "/var/lib/rancher/k3s/server/manifests"
)

var (
	ErrRequiredProperty = errors.New("property is required")
	ErrOutputOnly       = errors.New("property is only for output")
)

// TODO: generate schema from this struct
type Cluster struct {
	MasterNodes   []Node               `json:"masterNodes,omitempty"`
	Agents        []Node               `json:"agents,omitempty"`
	KubeConfig    string               `json:"kubeconfig,omitempty"`
	VersionConfig VersionConfiguration `json:"versionConfig,omitempty"`
}

// TODO: user and privatekey in provider as defaults (like region in openstack)
type Node struct {
	Host       string `json:"host,omitempty"`
	User       string `json:"user,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

// TODO: implement (put this into args section for master server?)
type Master struct {
	Node
	DisabledComponents []string // https://rancher.com/docs/k3s/latest/en/installation/install-options/server-config/#kubernetes-components
	FlannelBackend     string   //--flannel-backend
}

// VersionConfiguration resembles a K3s version. This can either be a release channel or a static version.
// If both are set the defined version will be preferred.
// Available channels can be found at: https://rancher.com/docs/k3s/latest/en/upgrades/basic/#release-channels
// An autoupdate configuration will automatically be added.
// For more information look here: https://rancher.com/docs/k3s/latest/en/upgrades/automated/
// If none is set stable channel will be used.
type VersionConfiguration struct {
	Channel string `json:"channel,omitempty"`
	Version string `json:"version,omitempty"`
}

func (v VersionConfiguration) EnvSetting() string {
	if v.Version != "" {
		return fmt.Sprintf("INSTALL_K3S_VERSION='%s'", v.Version)
	}

	channel := "stable"
	if v.Channel != "" {
		channel = v.Channel
	}

	return fmt.Sprintf("INSTALL_K3S_CHANNEL='%s'", channel)
}

func (v VersionConfiguration) YAMLValue() string {
	if v.Version != "" {
		return fmt.Sprintf("version: '%s'", v.Version)
	}

	channel := "stable"
	if v.Channel != "" {
		channel = v.Channel
	}

	return fmt.Sprintf("channel: '%s/%s'", channelURL, channel)
}

func MakeOrUpdateCluster(name string, cluster *Cluster) error {
	if err := cluster.Validate(); err != nil {
		return err
	}

	sudoPrefix := ""
	if useSudo {
		sudoPrefix = "sudo "
	}

	kubeconfig, err := setupNode(cluster.MasterNodes[0], cluster.VersionConfig, sudoPrefix)
	if err != nil {
		return err
	}

	if err := setupAutoUpdate(kubeconfig, cluster.VersionConfig); err != nil {
		return err
	}

	cluster.KubeConfig = kubeconfig

	return nil
}

func setupAutoUpdate(kubeconfig string, versionConfig VersionConfiguration) error {
	k8sClient, err := newK8sClient(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "failed to create client for kubernetes cluster")
	}

	manifestBuffer := &bytes.Buffer{}
	if err := upgradePlanManifestTemplate.Execute(manifestBuffer, versionConfig.YAMLValue()); err != nil {
		return errors.Wrapf(err, "failed to create upgradeplan")
	}

	for _, content := range [][]byte{systemUpgradeControllerManifest, manifestBuffer.Bytes()} {
		if err := k8sClient.CreateOrUpdateFromFile(context.Background(), content); err != nil {
			return errors.Wrap(err, "failed to apply k3s autoupdate objects")
		}
	}

	return nil
}

// NOTE: kept in case it's needed for containerd
func scpCopyManifests(sftpClient *sftp.Client, sshClient *sshexec.Client, fileReader io.Reader, fileName, sudoPrefix string) error {
	tmpFilePath := path.Join("/tmp", fileName)

	file, err := sftpClient.Create(tmpFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to create tmp file at %s", tmpFilePath)
	}

	if _, err := io.Copy(file, fileReader); err != nil {
		return errors.Wrapf(err, "failed to write to %s", tmpFilePath)
	}

	targetPath := path.Join(autoDeployManifestPath, fileName)
	mvCmd := &sshexec.Cmd{Command: fmt.Sprintf(
		"%smv %s %s", sudoPrefix, tmpFilePath, targetPath,
	)}
	if err := sshClient.Run(mvCmd); err != nil {
		return errors.Wrapf(err, "failed to move from %s to %s", tmpFilePath, targetPath)
	}

	return nil
}

// TODO: check if cluster is really working after intialization
// TODO: run kube-bench on k3s cluster
// TODO: implement hardening guide https://rancher.com/docs/k3s/latest/en/security/hardening_guide/
// TODO: add option to setup cilium as CNI
// https://docs.cilium.io/en/v1.9/gettingstarted/k3s/
// should be enough to enable ebf filesystem and disable the cni backend and then install cilium with kubernetes provider
// -> enableEbpf option?
// TODO: add option setup gVisor with containerd
// -> https://rancher.com/docs/k3s/latest/en/advanced/#configuring-containerd
// -> probably restart needed: sudo systemctl restart k3s
// -> maybe problems with https://github.com/k3s-io/k3s/issues/3378
func setupNode(node Node, versionConfig VersionConfiguration, sudoPrefix string) (string, error) {
	env := []string{
		versionConfig.EnvSetting(),
		fmt.Sprintf(`INSTALL_K3S_EXEC='server --tls-san "%s"'`, node.Host),
	}

	installK3scommand := fmt.Sprintf("%s | %s sh -\n", getScript, strings.Join(env, " "))

	getConfigcommand := fmt.Sprintf(sudoPrefix + "cat /etc/rancher/k3s/k3s.yaml")

	// execute commands

	// TODO: make port somehow configurable
	client, err := sshexec.NewClient(fmt.Sprintf("%s:22", node.Host), node.User, []byte(node.PrivateKey))
	if err != nil {
		return "", errors.Wrapf(err, "failed to establish ssh connection to %s", node.Host)
	}

	stdouterr := &bytes.Buffer{}
	err = client.Run(&sshexec.Cmd{
		Command: installK3scommand,
		Stderr:  stdouterr,
		Stdout:  stdouterr,
	})
	if err != nil {
		return "", errors.Wrap(errors.Wrap(err, stdouterr.String()), node.Host)
	}

	stdout := &bytes.Buffer{}
	err = client.Run(&sshexec.Cmd{
		Command: getConfigcommand,
		Stdout:  stdout,
	})
	if err != nil {
		return "", errors.Wrap(errors.Wrap(err, "failed to get kubeconfig"), node.Host)
	}

	kubeconfigReplacer := strings.NewReplacer(
		"127.0.0.1", node.Host,
		"localhost", node.Host,
	)

	return kubeconfigReplacer.Replace(stdout.String()), nil
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
			return errors.Wrapf(err, "failed to uninstall master %s", n.Host)
		}
	}

	for _, n := range cluster.Agents {
		if err := executeOnNode(n, "/usr/local/bin/k3s-agent-uninstall.sh"); err != nil {
			return errors.Wrapf(err, "failed to uninstall agent %s", n.Host)
		}
	}

	return nil
}

func executeOnNode(node Node, commands ...string) error {
	client, err := sshexec.NewClient(fmt.Sprintf("%s:22", node.Host), node.User, []byte(node.PrivateKey))
	if err != nil {
		return errors.Wrapf(err, "failed to establish ssh connection to %s", node.Host)
	}

	for _, command := range commands {
		stderr := &bytes.Buffer{}
		err = client.Run(&sshexec.Cmd{
			Command: command,
			Stderr:  stderr,
		})
		if err != nil {
			return errors.Wrap(err, stderr.String())
		}
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
