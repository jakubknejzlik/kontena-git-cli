package kontena

import (
	"fmt"
	"log"

	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/urfave/cli"

	"github.com/jakubknejzlik/kontena-git-cli/utils"
)

// CertificateInstall ...
func (c *Client) CertificateInstall(cert model.Certificate) error {
	utils.Log("installing certificate", cert.Domain)
	if cert.Bundle != "" {
		return c.DeployCertificate(cert, cert.Bundle)
	}

	if cert.Bundle == "" && cert.LetsEncrypt {
		return c.issueLECertificate(cert)
	}

	return cli.NewExitError(fmt.Sprintf(`certificate %s is not marked as letsencrypt and doesn't contain bundle`, cert.Domain), 1)
}

func (c *Client) issueLECertificate(cert model.Certificate) error {
	service := model.KontenaService{
		Environment: []string{
			"KONTENA_LB_VIRTUAL_HOSTS=" + cert.Domain,
			"KONTENA_LB_VIRTUAL_PATH=/.well-known/acme-challenge",
		},
		Links: []string{
			"core/internet_lb",
		},
		Image: "jakubknejzlik/acme.sh-nginx",
	}
	serviceName := "acme-challenge"

	c.ServiceRemove(serviceName)

	if err := c.ServiceCreate(serviceName, service); err != nil {
		return err
	}

	if err := c.ServiceDeploy(serviceName); err != nil {
		return err
	}

	utils.Log("issuing certificate")
	issueCmd := fmt.Sprintf(`/issue.sh %s`, cert.Domain)
	if data, err := c.ServiceExec(serviceName, issueCmd); err != nil {
		log.Println(err, string(data))
		// return err
	}

	utils.Log("fetching certificate")
	loadCertCmd := fmt.Sprintf(`cat /root/.acme.sh/%s/fullchain.cer /root/.acme.sh/%s/%s.key`, cert.Domain, cert.Domain, cert.Domain)
	if data, err := c.ServiceExec(serviceName, loadCertCmd); err == nil {
		c.DeployCertificate(cert, string(data))
	} else {
		return err
	}

	return c.removeAcmeService()
}

// DeployCertificate ...
func (c *Client) DeployCertificate(cert model.Certificate, bundle string) error {
	utils.Log("writing certificate", cert.SecretName())
	return c.WriteSecret(cert.SecretName(), bundle)
}

func (c *Client) removeAcmeService() error {
	utils.Log("removing acme-challenge service")
	removeServiceCmd := `kontena service remove --force acme-challenge`
	return utils.RunInteractive(removeServiceCmd)
}

// CurrentCertificateSecrets ...
func (c *Client) CurrentCertificateSecrets() ([]string, error) {
	certs := []string{}
	secrets, secretsErr := c.GetSecrets()

	if secretsErr != nil {
		return certs, secretsErr
	}

	for _, secretName := range secrets {
		if model.IsCertificateName(secretName) {
			certs = append(certs, secretName)
		}
	}

	return certs, nil
}

// kontena service create -e "KONTENA_LB_VIRTUAL_HOSTS=www.knejzlik.cz" -e "KONTENA_LB_VIRTUAL_PATH=/.well-known/acme-challenge" -l core/internet_lb acme-challenge jakubknejzlik/acme.sh-nginx
// kontena service deploy acme-challenge
// kontena service exec -it acme-challenge /issue.sh www.knejzlik.cz
// kontena service rm --force acme-challenge
//
//
// kontena service exec acme-challenge acme.sh --install-cert -d www.knejzlik.cz --key-file /key.pem --fullchain-file /cert.pem --reloadcmd "nginx -s reload" --debug
//
//
//
//
// kontena service create --cmd "daemon" test-acmesh neilpang/acme.sh
// kontena service deploy test-acmesh
// kontena service exec test-acmesh acme.sh --register-account
//
// kontena service exec test-acmesh acme.sh --issue -d www.knejzlik.cz  --stateless
// kontena service rm --force acme-challenge
// kontena service rm --force test-acmesh
//
// /usr/share/nginx/html
//
//
// kontena service exec test-acmesh nginx
// kontena service exec test-acmesh curl localhost
// kontena service create -e "STRING=aUKoszVA-SOvxjMWCzyWpEMidc1-rvRe0DWI_kXw52E" -e "KONTENA_LB_VIRTUAL_HOSTS=www.knejzlik.cz" -e "KONTENA_LB_VIRTUAL_PATH=/.well-known/acme-challenge" -l core/internet_lb --cmd "acme.sh --issue -d www.knejzlik.cz --stateless --debug" test-acmesh neilpang/acme.sh
//
// //aUKoszVA-SOvxjMWCzyWpEMidc1-rvRe0DWI_kXw52E
//
// kontena service create -e "STRING=MUGnf6WDNzXPgU_dkkFFhIFtUNhVAP31I09plw1pHrg" -e "KONTENA_LB_VIRTUAL_HOSTS=www.knejzlik.cz" -e "KONTENA_LB_VIRTUAL_PATH=/.well-known/acme-challenge" -l core/internet_lb --cmd "acme.sh --issue -d www.knejzlik.cz --stateless --debug" test-acmesh neilpang/acme.sh
// kontena service create --cmd "acme.sh --issue -d www.knejzlik.cz --stateless --debug" test-acmesh neilpang/acme.sh
// kontena service create -e "STRING=MUGnf6WDNzXPgU_dkkFFhIFtUNhVAP31I09plw1pHrg" -e "KONTENA_LB_VIRTUAL_HOSTS=www.knejzlik.cz" -e "KONTENA_LB_VIRTUAL_PATH=/.well-known/acme-challenge" --deploy-strategy daemon -l core/internet_lb test-nginx-string jakubknejzlik/nginx-string
// kontena service deploy test-acmesh && \
// kontena service logs -t test-acmesh
// kontena service rm --force test-acmesh
