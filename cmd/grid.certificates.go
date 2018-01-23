package cmd

import (
	"time"

	"github.com/inloop/goclitools"
	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"

	"github.com/urfave/cli"
)

func CertificatesCommand() cli.Command {
	return cli.Command{
		Name: "certificates",
		Subcommands: []cli.Command{
			installCertificatesCommand(),
			clearExpiredCertificatesCommand(),
		},
	}
}

func installCertificatesCommand() cli.Command {
	return cli.Command{
		Name: "install",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Install certificates")
			client := kontena.Client{}

			currentCertificateSecretsMap := map[string]bool{}
			currentCertificateSecrets, err := client.CurrentCertificateSecrets()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			for _, cert := range currentCertificateSecrets {
				currentCertificateSecretsMap[cert.Name] = true
			}

			certificates, err := model.CertificateLoadLocals()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			goclitools.Log("registered certificates:", len(currentCertificateSecrets), ", local certificates:", len(certificates))

			for _, certificate := range certificates {
				if currentCertificateSecretsMap[certificate.SecretName()] == false {
					if err := client.CertificateInstall(certificate); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
			}

			return nil
		},
	}
}

func clearExpiredCertificatesCommand() cli.Command {
	return cli.Command{
		Name: "clear",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Clear expired certificates")
			client := kontena.Client{}

			currentCertificateSecretsMap := map[string]bool{}
			currentCertificateSecrets, err := client.CurrentCertificateSecrets()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			for _, cert := range currentCertificateSecrets {
				currentCertificateSecretsMap[cert.Name] = true
			}

			secrets, err := client.SecretList()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			date := time.Now().AddDate(0, 0, 90-7)
			for _, secret := range secrets {
				if secret.IsCertificate() {
					if !secret.CreatedAt.Before(date) {
						goclitools.Log("removing certificate ", secret.Name, "; created:", secret.CreatedAt)
						if err := client.SecretRemove(secret.Name); err != nil {
							return cli.NewExitError(err, 1)
						}
					}
				}
			}

			goclitools.Log("all expiring certificates cleared")

			return nil
		},
	}
}
