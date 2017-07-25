package model

import (
	"io/ioutil"
	"os"
	"strings"

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
	ContainerName string            `yaml:"container_name,omitempty"`
	Instances     int               `yaml:"instances,omitempty"`
	Image         string            `yaml:"image,omitempty"`
	Command       string            `yaml:"command,omitempty"`
	Entrypoint    string            `yaml:"entrypoint,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Labels        map[string]string `yaml:"labels,omitempty"`
	Environment   []string          `yaml:"environment,omitempty"`
	Links         []string          `yaml:"links,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Secrets       []KontenaSecret   `yaml:"secrets,omitempty"`
	Deploy        KontenaDeploy     `yaml:"deploy,omitempty"`
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
					Name:   strings.Replace(secret.Name, "INLOOP_LB", "KONTENA_LB", 1),
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
