package k3s

import (
	_ "embed"
	"text/template"
)

//go:generate curl -sSfLO https://github.com/rancher/system-upgrade-controller/releases/latest/download/system-upgrade-controller.yaml

var (
	//go:embed system-upgrade-controller.yaml
	systemUpgradeControllerManifest []byte

	//go:embed containerd.toml.tmpl
	containerdConfigTemplate []byte

	//go:embed gvisor_install.sh
	gvisorInstall []byte

	gvisorUninstall = "rm -rf /usr/local/bin/runsc /usr/local/bin/containerd-shim-runsc-v1"

	//go:embed runtimeclass.yaml
	gvisorRuntimeClass []byte

	//go:embed cilium_install.sh
	ciliumInstall []byte

	//go:embed upgradeplan.yaml.tmpl
	upgradePlanManifestTemplateString string

	upgradePlanManifestTemplate = template.Must(template.New("").Parse(upgradePlanManifestTemplateString))
)
