package model

import (
	"io/ioutil"
	"os"

	"github.com/jakubknejzlik/kontena-git-cli/utils"

	yaml2 "gopkg.in/yaml.v2"
)

// ComposeSecret ...
type ComposeSecret struct {
	Secret string `yaml:"secret,omitempty"`
	Name   string `yaml:"name,omitempty"`
	Type   string `yaml:"type,omitempty"`
}

// ComposeDeploy ...
type ComposeDeploy struct {
	Strategy string `yaml:"strategy,omitempty"`
}

// ComposeService ...
type ComposeService struct {
	ContainerName string            `yaml:"container_name,omitempty"`
	Instances     string            `yaml:"instances,omitempty"` // this is not in standard docker-compose.yml !!!
	Image         string            `yaml:"image,omitempty"`
	Command       string            `yaml:"command,omitempty"`
	Entrypoint    string            `yaml:"entrypoint,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Labels        map[string]string `yaml:"labels,omitempty"`
	Environment   []string          `yaml:"environment,omitempty"`
	Links         []string          `yaml:"links,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Secrets       []ComposeSecret   `yaml:"secrets,omitempty"`
	Deploy        ComposeDeploy     `yaml:"deploy,omitempty"`
}

// Compose ...
type Compose struct {
	Stack    string                    `yaml:"stack,omitempty"` // this is not in standard docker-compose.yml !!!
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services,omitempty"`
}

// ComposeLoad ...
func ComposeLoad(path string) (Compose, error) {
	var dc Compose
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return dc, err
	}
	yaml2.Unmarshal(data, &dc)
	return dc, nil
}

// ExportTemporary ...
func (c *Compose) ExportTemporary() (string, error) {
	var path string
	tmp, err := ioutil.TempFile(os.TempDir(), "compose")
	if err != nil {
		return path, err
	}

	path = tmp.Name()

	data, marshalError := yaml2.Marshal(c)
	if marshalError != nil {
		return path, marshalError
	}

	utils.LogSection("exported compose file", string(data))
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return path, err
	}

	return path, nil
}
