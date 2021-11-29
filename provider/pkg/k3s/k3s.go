package k3s

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	getScript              = "curl -sfL https://get.k3s.io"
	useSudo                = true
	channelURL             = "https://update.k3s.io/v1-release/channels"
	autoDeployManifestPath = "/var/lib/rancher/k3s/server/manifests"
	containerdTemplatePath = "/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl"
)

var (
	ErrRequiredProperty = errors.New("property is required")
	ErrOutputOnly       = errors.New("property is only for output")
)

type Cluster struct {
	MasterNodes   []Node               `json:"masterNodes,omitempty"`
	Agents        []Node               `json:"agents,omitempty"`
	KubeConfig    string               `json:"kubeconfig,omitempty"`
	VersionConfig VersionConfiguration `json:"versionConfig,omitempty"`
	CNIConfig     CNIConfig            `json:"cniConfig,omitempty"`
}

type Node struct {
	Host       string `json:"host,omitempty"`
	User       string `json:"user,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	// Args define CLI arguments for k3s server or k3s agent respectively.
	// The passed args won't be validated and just passed to the installation instructions of the node.
	// An example value for the master node would look like []string{"--disable=traefik"}.
	Args      []string  `json:"args,omitempty"`
	CRIConfig CRIConfig `json:"criConfig,omitempty"`
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

type CRIConfig struct {
	EnableGVisor bool `json:"enableGVisor,omitempty"`
}

type CNIConfig struct {
	Provider string `json:"provider,omitempty"`
}

func MakeOrUpdateCluster(name string, cluster *Cluster) error {
	if err := cluster.Validate(); err != nil {
		return err
	}

	kubeconfig, err := setupNode(cluster.MasterNodes[0], cluster.VersionConfig)
	if err != nil {
		return err
	}

	if err := setupAutoUpdate(kubeconfig, cluster.VersionConfig); err != nil {
		return err
	}

	if cluster.MasterNodes[0].CRIConfig.EnableGVisor {
		if err := setupGVisor(cluster.MasterNodes[0], kubeconfig); err != nil {
			return err
		}
	} else {
		if err := uninstallGVisor(cluster.MasterNodes[0], kubeconfig); err != nil {
			return err
		}
	}

	cluster.KubeConfig = kubeconfig

	return nil
}

func setupGVisor(node Node, kubeconfig string) error {
	remoteExecutor, err := NewExecutorForNode(node, useSudo)
	if err != nil {
		return err
	}

	if err := remoteExecutor.CopyFile(bytes.NewReader(containerdConfigTemplate), containerdTemplatePath); err != nil {
		return err
	}

	if err := remoteExecutor.ExecuteScript(gvisorInstall); err != nil {
		return err
	}

	if _, err := remoteExecutor.SudoCombinedOutput("systemctl restart k3s.service"); err != nil {
		return err
	}

	k8sClient, err := newK8sClient(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "failed to create client for kubernetes cluster")
	}

	if err := k8sClient.CreateOrUpdateFromFile(context.TODO(), gvisorRuntimeClass); err != nil {
		return err
	}

	return nil
}

func uninstallGVisor(node Node, kubeconfig string) error {
	k8sClient, err := newK8sClient(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "failed to create client for kubernetes cluster")
	}

	if err := k8sClient.DeleteIfExistsFromFile(context.TODO(), gvisorRuntimeClass); err != nil {
		return err
	}

	remoteExecutor, err := NewExecutorForNode(node, useSudo)
	if err != nil {
		return err
	}

	for _, cmd := range []string{
		gvisorUninstall, fmt.Sprintf("rm %s", containerdTemplatePath), "systemctl restart k3s.service",
	} {
		if _, err := remoteExecutor.SudoCombinedOutput(cmd); err != nil {
			return err
		}
	}

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

func setupNode(node Node, versionConfig VersionConfiguration) (string, error) {
	remoteExecutor, err := NewExecutorForNode(node, useSudo)
	if err != nil {
		return "", err
	}

	env := []string{
		versionConfig.EnvSetting(),
		fmt.Sprintf(`INSTALL_K3S_EXEC='server --tls-san="%s" %s'`, node.Host, strings.Join(node.Args, " ")),
	}

	installK3scommand := fmt.Sprintf("%s | %s sh -\n", getScript, strings.Join(env, " "))

	_, err = remoteExecutor.CombinedOutput(installK3scommand)
	if err != nil {
		return "", err
	}

	kubeconfig, err := remoteExecutor.SudoOutput("cat /etc/rancher/k3s/k3s.yaml")
	if err != nil {
		return "", err
	}

	kubeconfigReplacer := strings.NewReplacer(
		"127.0.0.1", node.Host,
		"localhost", node.Host,
	)

	return kubeconfigReplacer.Replace(kubeconfig), nil
}

// TODO: add version validation (only last 3 versions are supported because of containerd.toml)
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
		if err := removeNode(n); err != nil {
			return err
		}
	}

	// for _, n := range cluster.Agents {
	// 	if err := executeScriptIfExistsOnNode(
	// 		n, "sh /usr/local/bin/k3s-agent-uninstall.sh", gvisorUninstall,
	// 	); err != nil {
	// 		return errors.Wrapf(err, "failed to uninstall agent")
	// 	}
	// }

	return nil
}

func removeNode(node Node) error {
	remoteExecutor, err := NewExecutorForNode(node, useSudo)
	if err != nil {
		return err
	}

	masterUninstall := "/usr/local/bin/k3s-uninstall.sh"
	uninstallCmds := []string{gvisorUninstall}

	if _, err := remoteExecutor.FileHandler.Stat(masterUninstall); err == nil {
		uninstallCmds = append(uninstallCmds, fmt.Sprintf("sh %s", masterUninstall))
	}

	for _, cmd := range uninstallCmds {
		if _, err := remoteExecutor.SudoCombinedOutput(cmd); err != nil {
			return err
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
