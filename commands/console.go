package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// ConsoleCmd ...
	ConsoleCmd = &cobra.Command{
		Use:   "console",
		Short: "Open an interactive console inside a component.",
		Long:  ``,
		Run:   consoleFn,
	}
	user string
)

func init() {
	ConsoleCmd.Flags().StringVarP(&user, "user", "u", "", "user you would like to console in as")
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {
	if user != "" {
		registry.Set("console_user", user)
	}
	envModel, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(envModel, args, 2)

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will os.Exit(1) with an error message
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the component you wish to console into:

ex: nanobox console local web.site

`, len(args))
		return
	}

	switch location {
	case "local":
		appModel, _ := models.FindAppBySlug(config.EnvID(), name)
		if appModel.Status != "up" {
			fmt.Println("unable to continue until the app is up")
			return
		}

		componentModel, _ := models.FindComponentBySlug(config.EnvID()+"_"+name, args[0])
		// todo: determine ways this errors and handle it, use util.Err for suggestions
		// componentModel, err := models.FindComponentBySlug(config.EnvID()+"_"+name, args[0])
		// if err != nil {
		// 	display.CommandErr(err)
		// 	return
		// }

		display.CommandErr(env.Console(componentModel, console.ConsoleConfig{}))

	case "production":

		consoleConfig := processors.ConsoleConfig{
			App:  name,
			Host: args[0],
		}

		// set the meta arguments to be used in the processor and run the processor
		display.CommandErr(processors.Console(envModel, consoleConfig))

	}
}
