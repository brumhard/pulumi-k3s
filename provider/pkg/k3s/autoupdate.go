package k3s

import (
	_ "embed"
	"text/template"
)

//go:generate curl -sSfLO https://github.com/rancher/system-upgrade-controller/releases/latest/download/system-upgrade-controller.yaml

var (
	//go:embed system-upgrade-controller.yaml
	systemUpgradeControllerManifest []byte

	//go:embed upgradeplan.yaml.tmpl
	upgradePlanManifestTemplateString string

	upgradePlanManifestTemplate = template.Must(template.New("").Parse(upgradePlanManifestTemplateString))
)
