package cmd

import (
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
		},
	}
}

func installCertificatesCommand() cli.Command {
	return cli.Command{
		Name: "install",
		Action: func(c *cli.Context) error {
			goclitools.LogSection("Install certificates")
			client := kontena.Client{}

			currentCertificatesMap := map[string]bool{}
			currentCertificates, err := client.CurrentCertificates()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			for _, cert := range currentCertificates {
				currentCertificatesMap[cert] = true
			}

			certificates, err := model.CertificateLoadLocals()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			goclitools.Log("registered certificates:", len(currentCertificates), ", local certificates:", len(certificates))

			for _, certificate := range certificates {
				if currentCertificatesMap[certificate.Domain] == false {
					if err := client.CertificateInstall(certificate); err != nil {
						return cli.NewExitError(err, 1)
					}
				}
			}

			return nil
		},
	}
}
