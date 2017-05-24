package model

import (
	"io/ioutil"

	yaml2 "gopkg.in/yaml.v2"
)

// Registry ...
type Registry struct {
	Name     string `yaml:"name,omitempty"`
	User     string `yaml:"username,omitempty"`
	Email    string `yaml:"email,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// RegistriesLoad ...
func RegistriesLoad(path string) ([]Registry, error) {
	var result []Registry

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}

	var regs map[string]Registry
	yaml2.Unmarshal(data, &regs)

	for name, registry := range regs {
		registry.Name = name
		result = append(result, registry)
	}

	return result, nil
}
