package model

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/jakubknejzlik/kontena-git-cli/utils"

	yaml2 "gopkg.in/yaml.v2"
)

// Certificate ...
type Certificate struct {
	Domain           string   `yaml:"domain,omitempty"`
	Bundle           string   `yaml:"bundle,omitempty"`
	Type             string   `yaml:"type,omitempty"`
	AlternativeNames []string `yaml:"alternative_names,omitempty"`
}

// CertificateLoadLocals ...
func CertificateLoadLocals() (map[string]Certificate, error) {
	var certs map[string]Certificate

	if utils.FileExists("certificates.yml") {
		data, err := ioutil.ReadFile("certificates.yml")
		if err != nil {
			return certs, err
		}

		yaml2.Unmarshal(data, &certs)
		for domain, cert := range certs {
			cert.Domain = domain
			certs[domain] = cert
		}
	} else {
		certs = map[string]Certificate{}
	}

	certsFiles, _ := ioutil.ReadDir("./certificates")
	for _, certFile := range certsFiles {
		certName := certFile.Name()
		data, dataErr := ioutil.ReadFile(path.Join("./certificates", certName))
		if dataErr != nil {
			return certs, dataErr
		}
		cert := Certificate{
			Domain: certName,
			Bundle: string(data),
		}
		certs[certName] = cert
	}

	return certs, nil
}

// AllDomains ...
func (c Certificate) AllDomains() []string {
	return append([]string{c.Domain}, c.AlternativeNames...)
}

// SecretName ...
func (c Certificate) SecretName() string {
	rg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	name := rg.ReplaceAllString(c.Domain, "_")

	return "core_SSL_CERTIFICATE_" + name + "_BUNDLE"
}

// Description ...
func (c Certificate) Description() string {
	return fmt.Sprintf("%s (%s)", c.Domain, strings.Join(c.AlternativeNames, ", "))
}
