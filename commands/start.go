package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/processors/provider/bridge"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/service"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the Nanobox virtual machine.",
		Long:  ``,
		Run:   startFn,
	}
)

func init() {
	steps.Build("start", startCheck, startFn)
}

// startFn ...
func startFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(processors.Start())
}

func startCheck() bool {
	bridgeReady := true
	if provider.BridgeRequired() {
		bridgeReady = bridge.Connected()
	}
	return provider.IsReady() && service.Running("nanobox-server") && bridgeReady
}
