package model

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/jakubknejzlik/kontena-git-cli/utils"

	yaml2 "gopkg.in/yaml.v2"
)

// Certificate ...
type Certificate struct {
	Domain      string `yaml:"domain,omitempty"`
	Bundle      string `yaml:"bundle,omitempty"`
	LetsEncrypt bool   `yaml:"letsencrypt,omitempty"`
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

// SecretName ...
func (c Certificate) SecretName() string {
	rg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	name := rg.ReplaceAllString(c.Domain, "_")

	return "core_SSL_CERTIFICATE_" + name + "_BUNDLE"
}

// IsCertificateName ...
func IsCertificateName(name string) bool {
	re := regexp.MustCompile(`core_SSL_CERTIFICATE_[a-zA-Z0-9_]+_BUNDLE`)
	return re.Match([]byte(name))
}
