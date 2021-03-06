package commands

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
)

type Api struct {
	ui           terminal.UI
	endpointRepo api.EndpointRepository
	config       configuration.Reader
}

type ApiEndpointSetter interface {
	SetApiEndpoint(endpoint string)
}

func NewApi(ui terminal.UI, config configuration.Reader, endpointRepo api.EndpointRepository) (cmd Api) {
	cmd.ui = ui
	cmd.config = config
	cmd.endpointRepo = endpointRepo
	return
}

func (cmd Api) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	return
}

func (cmd Api) Run(c *cli.Context) {
	if len(c.Args()) == 0 {
		if cmd.config.ApiEndpoint() == "" {
			cmd.ui.Say(fmt.Sprintf("No api endpoint set. Use '%s' to set an endpoint", terminal.CommandColor(cf.Name()+" api")))
		} else {
			cmd.ui.Say(
				"API endpoint: %s (API version: %s)",
				terminal.EntityNameColor(cmd.config.ApiEndpoint()),
				terminal.EntityNameColor(cmd.config.ApiVersion()),
			)
		}
		return
	}

	cmd.SetApiEndpoint(c.Args()[0])
}

func (cmd Api) SetApiEndpoint(endpoint string) {
	if strings.HasSuffix(endpoint, "/") {
		endpoint = strings.TrimSuffix(endpoint, "/")
	}

	cmd.ui.Say("Setting api endpoint to %s...", terminal.EntityNameColor(endpoint))

	endpoint, apiErr := cmd.endpointRepo.UpdateEndpoint(endpoint)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if !strings.HasPrefix(endpoint, "https://") {
		cmd.ui.Say(terminal.WarningColor("Warning: Insecure http API endpoint detected: secure https API endpoints are recommended\n"))
	}

	cmd.ui.ShowConfiguration(cmd.config)
}
