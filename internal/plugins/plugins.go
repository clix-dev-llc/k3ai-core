package plugins

import (
	"github.com/kf5i/k3ai-core/internal/shared"
	"github.com/pkg/errors"
)

const (
	commandFile = "file"
	// CommandKustomize is the kustomize command
	CommandKustomize = "kustomize"
	// DefaultPluginFileName is the default plugin name
	// each plugin must contain this file else it will be ignored
	DefaultPluginFileName = "plugin.yaml"
	// DefaultGroupFileName is the default group name
	// each group must contain this file else it will be ignored
	DefaultGroupFileName = "group.yaml"
)

// YamlType is the specification for YamlType segment of the Plugin
type YamlType struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type,omitempty"`
}

//PostInstall to execute after the scripts
type PostInstall struct {
	Command string `yaml:"command,omitempty"`
}

//Plugin is the specification of each k3ai plugin
type Plugin struct {
	Namespace         string      `yaml:"namespace,omitempty"`
	Labels            []string    `yaml:",flow"`
	PluginName        string      `yaml:"plugin-name"`
	PluginDescription string      `yaml:"plugin-description"`
	Yaml              []YamlType  `yaml:"yaml,flow"`
	Bash              []string    `yaml:"bash,flow"`
	Helm              []string    `yaml:"helm,flow"`
	PostInstall       PostInstall `yaml:"post-install"`
}

// Plugins list of plugins
type Plugins struct {
	Items []Plugin `yaml:"items"`
}

// Encode fetches the Plugin
func (ps *Plugin) Encode(URL string) error {
	err := encode(URL, ps)
	if err != nil {
		return err
	}
	mergeWithDefault(ps)
	return nil
}

// List fetch the plugin list
func (pls *Plugins) List(URL string) error {
	gHubContents, err := GithubContentList(URL)
	if err != nil {
		return err
	}
	for _, gHubContent := range gHubContents {
		var ps Plugin
		err := ps.Encode(shared.NormalizeURL(URL, gHubContent.Name) + DefaultPluginFileName)
		if err != nil {
			return err
		}
		pls.Items = append(pls.Items, ps)
	}
	return nil
}

// validate checks for any errors in the Plugin
func (ps *Plugin) validate() error {
	if ps.Namespace == "" {
		return errors.New("namespace value must be 'default' or another value")
	}
	for _, yamlType := range ps.Yaml {
		if yamlType.Type != CommandKustomize && yamlType.Type != commandFile {
			return errors.New("type must be file or kustomize")
		}
	}
	return nil
}

func mergeWithDefault(ps *Plugin) {
	ps.Namespace = shared.GetDefaultIfEmpty(ps.Namespace, "default")
	for i, yamlTypeItem := range ps.Yaml {
		yamlType := shared.GetDefaultIfEmpty(yamlTypeItem.Type, "file")
		ps.Yaml[i] = YamlType{Type: yamlType, URL: yamlTypeItem.URL}
	}
}
