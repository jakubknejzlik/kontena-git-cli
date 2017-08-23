package cmd

import (
	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/jakubknejzlik/kontena-git-cli/utils"

	"github.com/urfave/cli"
)

func installCertificatesCommand() cli.Command {
	return cli.Command{
		Name: "certificates",
		Action: func(c *cli.Context) error {
			utils.LogSection("Certificates")
			client := kontena.Client{}

			currentCertificateSecretsMap := map[string]bool{}
			currentCertificateSecrets, err := client.CurrentCertificateSecrets()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			for _, cert := range currentCertificateSecrets {
				currentCertificateSecretsMap[cert] = true
			}

			certificates, err := model.CertificateLoadLocals()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			utils.Log("registered certificates:", len(currentCertificateSecrets), ", local certificates:", len(certificates))

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
