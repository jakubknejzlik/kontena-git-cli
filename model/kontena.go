package model

import (
	"io/ioutil"
	"os"

	"github.com/inloop/goclitools"
	yaml2 "gopkg.in/yaml.v2"
)

// KontenaStack ...
type KontenaStack struct {
	Name     string                    `yaml:"stack,omitempty"`
	Version  string                    `yaml:"version,omitempty"`
	Expose   string                    `yaml:"expose,omitempty"`
	Services map[string]KontenaService `yaml:"services"`
	Volumes  map[string]KontenaVolume  `yaml:"volumes"`
}

// KontenaSecret ...
type KontenaSecret struct {
	Secret string `yaml:"secret,omitempty"`
	Name   string `yaml:"name,omitempty"`
	Type   string `yaml:"type,omitempty"`
}

// KontenaService ...
type KontenaService struct {
	Image           string                        `yaml:"image,omitempty"`
	Instances       *int                          `yaml:"instances,omitempty"`
	Stateful        bool                          `yaml:"stateful,omitempty"`
	Command         string                        `yaml:"command,omitempty"`
	Volumes         []string                      `yaml:"volumes"`
	VolumesFrom     []string                      `yaml:"volumes_from"`
	Environment     []string                      `yaml:"environment"`
	EnvFile         string                        `yaml:"env_file,omitempty"`
	Links           []string                      `yaml:"links"`
	DependsOn       []string                      `yaml:"depends_on"`
	Ports           []string                      `yaml:"ports"`
	Affinity        []string                      `yaml:"affinity"`
	CapAdd          []string                      `yaml:"cap_add"`
	CapDrop         []string                      `yaml:"cap_drop"`
	CPUShares       *int                          `yaml:"cpu_shares"`
	MemLimit        string                        `yaml:"mem_limit,omitempty"`
	MemswapLimit    string                        `yaml:"memswap_limit,omitempty"`
	StopGracePeriod string                        `yaml:"stop_grace_period,omitempty"`
	NetworkMode     string                        `yaml:"network_mode,omitempty"`
	Pid             string                        `yaml:"pid,omitempty"`
	Privileged      bool                          `yaml:"privileged,omitempty"`
	User            string                        `yaml:"user,omitempty"`
	Secrets         []KontenaSecret               `yaml:"secrets"`
	Hooks           map[string]KontenaServiceHook `yaml:"hooks"`
	Extends         KontenaServiceExtends         `yaml:"extends,omitempty"`
	Deploy          KontenaServiceDeploy          `yaml:"deploy,omitempty"`
	Logging         KontenaServiceLogging         `yaml:"logging,omitempty"`
	HealthCheck     KontenaServiceHealthCheck     `yaml:"health_check,omitempty"`
}

// KontenaServiceExtends ...
type KontenaServiceExtends struct {
	File    string `yaml:"file,omitempty"`
	Service string `yaml:"service,omitempty"`
}

// KontenaServiceDeploy ...
type KontenaServiceDeploy struct {
	Strategy    string  `yaml:"strategy,omitempty"`
	WaitForPort int     `yaml:"wait_for_port,omitempty"`
	MinHealth   float32 `yaml:"min_health,omitempty"`
	Interval    string  `yaml:"interval,omitempty"`
}

// KontenaServiceHook ...
type KontenaServiceHook struct {
	Name     string `yaml:"name,omitempty"`
	Cmd      string `yaml:"cmd,omitempty"`
	Instance string `yaml:"instance,omitempty"`
	OneShot  bool   `yaml:"one_shot,omitempty"`
}

// KontenaServiceLogging ...
type KontenaServiceLogging struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

// KontenaServiceHealthCheck ...
type KontenaServiceHealthCheck struct {
	Protocol     string `yaml:"protocol,omitempty"`
	Port         int    `yaml:"port,omitempty"`
	Interval     int    `yaml:"interval,omitempty"`
	URI          string `yaml:"uri,omitempty"`
	InitialDelay int    `yaml:"initial_delay,omitempty"`
	Timeout      int    `yaml:"timeout,omitempty"`
}

// KontenaVolume ...
type KontenaVolume struct {
	External bool `yaml:"external,omitempty"`
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

	goclitools.LogSection("exported kontena file", string(data))
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return path, err
	}

	return path, nil
}
