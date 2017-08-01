package model

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"

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

	data, err := ioutil.ReadFile("certificates.yml")
	if err != nil {
		return certs, err
	}

	yaml2.Unmarshal(data, &certs)
	for domain, cert := range certs {
		cert.Domain = domain
		certs[domain] = cert
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
	return "core_SSL_CERTIFICATE_" + strings.Replace(c.Domain, ".", "_", -1) + "_BUNDLE"
}

// IsCertificateName ...
func IsCertificateName(name string) bool {
	re := regexp.MustCompile(`core_SSL_CERTIFICATE_[a-z0-9\*_]+_BUNDLE`)
	return re.Match([]byte(name))
}
