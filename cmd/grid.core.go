package cmd

import (
	"strings"

	"github.com/jakubknejzlik/kontena-git-cli/utils"

	"github.com/jakubknejzlik/kontena-git-cli/kontena"
	"github.com/jakubknejzlik/kontena-git-cli/model"
	"github.com/urfave/cli"
)

func installCoreCommand() cli.Command {
	return cli.Command{
		Name: "core",
		Action: func(c *cli.Context) error {
			utils.LogSection("Core")
			client := kontena.Client{}

			dc, err := model.KontenaLoad("kontena.yml")
			if err != nil {
				dc = defaultRootStack()
			}

			if dc.Services["internet_lb"].Image == "" {
				dc.Services["internet_lb"] = defaultLoadBalancer()
			}

			loadBalancer := dc.Services["internet_lb"]
			currentCertificates, certsErr := client.CurrentCertificateSecrets()
			if certsErr != nil {
				return cli.NewExitError(certsErr, 1)
			}

			for _, secret := range currentCertificates {
				s := model.KontenaSecret{
					Secret: strings.Replace(secret, "core_", "", 1),
					Name:   "SSL_CERTS",
					Type:   "env",
				}
				loadBalancer.Secrets = append(loadBalancer.Secrets, s)
			}
			dc.Services["internet_lb"] = loadBalancer

			if err := client.StackInstallOrUpgrade(dc); err != nil {
				return cli.NewExitError(err, 1)
			}

			if err := client.StackDeploy("core"); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func defaultRootStack() model.KontenaStack {
	return model.KontenaStack{
		Name:    "core",
		Version: "1.0.0",
		Services: map[string]model.KontenaService{
			"internet_lb": defaultLoadBalancer(),
		},
	}
}

func defaultLoadBalancer() model.KontenaService {
	return model.KontenaService{
		Image: "kontena/lb:latest",
		Ports: []string{"80:80", "443:443"},
		Deploy: model.KontenaDeploy{
			Strategy: "daemon",
		},
	}
}
