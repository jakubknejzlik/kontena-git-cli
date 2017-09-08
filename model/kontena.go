package model

import (
	"io/ioutil"
	"os"

	"github.com/jakubknejzlik/kontena-git-cli/utils"

	yaml2 "gopkg.in/yaml.v2"
)

// KontenaSecret ...
type KontenaSecret struct {
	Secret string `yaml:"secret,omitempty"`
	Name   string `yaml:"name,omitempty"`
	Type   string `yaml:"type,omitempty"`
}

// KontenaDeploy ...
type KontenaDeploy struct {
	Strategy string `yaml:"strategy,omitempty"`
}

// KontenaService ...
type KontenaService struct {
	Instances *int   `yaml:"instances,omitempty"`
	Image     string `yaml:"image,omitempty"`
	Command   string `yaml:"command,omitempty"`
	// Entrypoint    string            `yaml:"entrypoint,omitempty"`
	Volumes     []string              `yaml:"volumes"`
	Labels      map[string]string     `yaml:"labels"`
	Environment []string              `yaml:"environment"`
	Links       []string              `yaml:"links"`
	Ports       []string              `yaml:"ports"`
	Secrets     []KontenaSecret       `yaml:"secrets"`
	Deploy      KontenaDeploy         `yaml:"deploy,omitempty"`
	Logging     KontenaServiceLogging `yaml:"logging,omitempty"`
	Stateful    bool                  `yaml:"stateful,omitempty"`
}

// KontenaServiceLogging ...
type KontenaServiceLogging struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

// KontenaStack ...
type KontenaStack struct {
	Name     string                    `yaml:"stack,omitempty"`
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]KontenaService `yaml:"services,omitempty"`
}

// KontenaLoad ...
func KontenaLoad(path string) (KontenaStack, error) {
	var dc KontenaStack
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return dc, err
	}
	yaml2.Unmarshal(data, &dc)
	return dc, nil
}

// ExportTemporary ...
func (c KontenaStack) ExportTemporary(translateSecrets bool) (string, error) {
	var path string
	tmp, err := ioutil.TempFile(os.TempDir(), "kontena")
	if err != nil {
		return path, err
	}

	path = tmp.Name()
	stack := c

	if translateSecrets {
		newServices := map[string]KontenaService{}
		for i, service := range stack.Services {
			newSecrets := []KontenaSecret{}
			for _, secret := range service.Secrets {
				newSecrets = append(newSecrets, KontenaSecret{
					Secret: stack.Name + "_" + secret.Secret,
					Name:   secret.Name,
					Type:   secret.Type,
				})
			}
			service.Secrets = newSecrets
			newServices[i] = service
		}
		stack.Services = newServices
	}

	data, marshalError := yaml2.Marshal(stack)
	if marshalError != nil {
		return path, marshalError
	}

	utils.LogSection("exported kontena file", string(data))
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return path, err
	}

	return path, nil
}
