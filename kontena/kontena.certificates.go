package kontena

import (
	"fmt"
	"strings"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"
	"github.com/urfave/cli"
)

// CertificateInstall ...
func (c *Client) CertificateInstall(cert model.Certificate) error {
	return c.CertificateInstallInGrid(c.CurrentGrid().Name, cert)
}

// CertificateInstallInGrid ...
func (c *Client) CertificateInstallInGrid(grid string, cert model.Certificate) error {
	goclitools.Log("installing certificate", cert.Description(), "in grid", grid)
	if cert.Bundle != "" {
		return c.DeployCertificateInGrid(grid, cert, cert.Bundle)
	}

	if cert.Bundle == "" {
		return c.issueLECertificateInGrid(grid, cert)
	}

	return cli.NewExitError(fmt.Sprintf(`certificate %s is not marked as letsencrypt and doesn't contain bundle`, cert.Domain), 1)
}

// DeployCertificateInGrid ...
func (c *Client) DeployCertificateInGrid(grid string, cert model.Certificate, bundle string) error {
	goclitools.Log("writing certificate", cert.SecretName(), "grid", grid)
	return c.SecretWriteToGrid(grid, cert.SecretName(), bundle)
}

func (c *Client) issueLECertificateInGrid(grid string, cert model.Certificate) error {
	allDomains := cert.AllDomains()

	for _, domain := range allDomains {
		if err := authorizeLetsEncryptDomain(domain, cert.Type, grid); err != nil {
			return err
		}
	}

	return requestLetsEncryptDomains(allDomains, grid)
}

func authorizeLetsEncryptDomain(domain, authType, grid string) error {
	if authType == "" {
		authType = "http-01"
	}
	cmd := fmt.Sprintf("kontena certificate authorize --grid %s --type %s --linked-service core/internet_lb %s", grid, authType, domain)
	return goclitools.RunInteractive(cmd)
}

func requestLetsEncryptDomains(domains []string, grid string) error {
	cmd := fmt.Sprintf("kontena certificate request --grid %s %s", grid, strings.Join(domains, " "))
	return goclitools.RunInteractive(cmd)
}

// CurrentCertificates ...
func (c *Client) CurrentCertificates() ([]string, error) {
	var certs = []string{}
	cmd := fmt.Sprintf("kontena certificate ls -q")
	res, err := goclitools.Run(cmd)
	if err != nil {
		return certs, err
	}
	certs = utils.SplitString(string(res), "\n")
	return certs, nil
}
