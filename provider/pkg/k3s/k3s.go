package k3s

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/brumhard/pulumi-k3s/provider/pkg/sshexec"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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
	MasterNodes   []Node               `json:"masterNodes"`
	Agents        []Node               `json:"agents"`
	KubeConfig    string               `json:"kubeconfig"`
	VersionConfig VersionConfiguration `json:"versionConfig"`
}

// TODO: user and privatekey in provider as defaults (like region in openstack)
type Node struct {
	Host       string `json:"host"`
	User       string `json:"user"`
	PrivateKey string `json:"privateKey"`
}


// VersionConfiguration resembles a K3s version. This can either be a release channel or a static version.
// If both are set the defined version will be preferred.
// Available channels can be found at: https://rancher.com/docs/k3s/latest/en/upgrades/basic/#release-channels
// An autoupdate configuration will automatically be added.
// For more information look here: https://rancher.com/docs/k3s/latest/en/upgrades/automated/
// If none is set stable channel will be used.
type VersionConfiguration struct {
	Channel string `json:"channel"`
	Version string `json:"version"`
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

	if err := setupAutoUpdate(cluster.MasterNodes[0], cluster.VersionConfig, sudoPrefix); err != nil {
		return err
	}

	cluster.KubeConfig = kubeconfig

	return nil
}

func setupAutoUpdate(node Node, versionConfig VersionConfiguration, sudoPrefix string) error {
	signer, err := ssh.ParsePrivateKey([]byte(node.PrivateKey))
	if err != nil {
		return err
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", node.Host), &ssh.ClientConfig{
		User:            node.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return err
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}

	sshClient := sshexec.NewClientFromSSH(conn)

	tmpl, err := template.New("").Parse(upgradePlanManifestTemplate)
	if err != nil {
		return err
	}

	manifestBuffer := &bytes.Buffer{}
	if err := tmpl.Execute(manifestBuffer, versionConfig.YAMLValue()); err != nil {
		return err
	}

	filesMap := map[string]io.Reader{
		"system-upgrade-controller.yaml": bytes.NewBuffer(systemUpgradeControllerManifest),
		"upgradeplan.yaml":               manifestBuffer,
	}

	for name, content := range filesMap {
		if err := scpCopyManifests(sftpClient, sshClient, content, name, sudoPrefix); err != nil {
			return err
		}
	}

	return nil
}

func scpCopyManifests(sftpClient *sftp.Client, sshClient *sshexec.Client, fileReader io.Reader, fileName, sudoPrefix string) error {
	tmpFilePath := path.Join("/tmp", fileName)

	file, err := sftpClient.Create(tmpFilePath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, fileReader); err != nil {
		return err
	}

	mvCmd := &sshexec.Cmd{Command: fmt.Sprintf(
		"%smv %s %s", sudoPrefix, tmpFilePath, path.Join(autoDeployManifestPath, fileName),
	)}
	if err := sshClient.Run(mvCmd); err != nil {
		return err
	}

	return nil
}

// TODO: check if cluster is really working after intialization
// TODO: run kube-bench on k3s cluster
// TODO: add option to setup cilium as CNI
// https://docs.cilium.io/en/v1.9/gettingstarted/k3s/
// should be enough to enable ebf filesystem and disable the cni backend and then install cilium with kubernetes provider
// -> enableEbpf option?
// TODO: add option setup gVisor with containerd
// -> https://rancher.com/docs/k3s/latest/en/advanced/#configuring-containerd
// -> probably restart needed: sudo systemctl restart k3s
// -> maybe problems with https://github.com/k3s-io/k3s/issues/3378
// TODO: implement hardening guide https://rancher.com/docs/k3s/latest/en/security/hardening_guide/
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
		return "", err
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

func executeOnNode(node Node, commands ...string) error {
	client, err := sshexec.NewClient(fmt.Sprintf("%s:22", node.Host), node.User, []byte(node.PrivateKey))
	if err != nil {
		return err
	}

	for _, command := range commands {
		stderr := &bytes.Buffer{}
		err = client.Run(&sshexec.Cmd{
			Command: command,
			Stderr:  stderr,
		})
		if err != nil {
			return errors.Wrap(errors.Wrap(err, stderr.String()), node.Host)
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
