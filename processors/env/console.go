package env

import (
	"fmt"
	"os"
	"io"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/docker/docker/pkg/term"
	"github.com/jcelliott/lumber"


	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/provider"
)

// Console ...
func Console(componentModel *models.Component, consoleConfig ConsoleConfig) error {
	// set the default shell
	if consoleConfig.Shell == "" {
		consoleConfig.Shell = "bash"
	}

	// setup docker client
	if err := provider.Init(); err != nil {
		return err
	}

	// this is the default command to run in the container
	cmd := []string{"/bin/bash"}

	// check to see if there are any optional meta arguments that need to be handled
	switch {

	// if a current working directory (cwd) is provided then modify the command to
	// change into that directory before executing
	case consoleConfig.Cwd != "":
		cmd = append(cmd, "-c", fmt.Sprintf("cd %s; exec \"%s\"", consoleConfig.Cwd, consoleConfig.Shell))

	// if a command is provided then modify the command to exec that command after
	// running the base command
	case consoleConfig.Command != "":
		cmd = append(cmd, "-c", consoleConfig.Command)
	}

	// establish file descriptors for std streams
	stdInFD, isTerminal := term.GetFdInfo(os.Stdin)
	stdOutFD, _ := term.GetFdInfo(os.Stdout)

	// initiate a docker exec
	_, resp, err := docker.ExecStart(componentModel.ID, cmd, true, true, true, isTerminal)
	if err != nil {
		lumber.Error("dockerexecerror: %s", err)
		return err
	}
	defer resp.Conn.Close()

	// if we are using a term, lets upgrade it to RawMode
	if isTerminal {
		oldInState, err := term.SetRawTerminal(stdInFD)
		if err == nil {
			defer term.RestoreTerminal(stdInFD, oldInState)
		}

		oldOutState, err := term.SetRawTerminalOutput(stdOutFD)
		if err == nil {
			defer term.RestoreTerminal(stdOutFD, oldOutState)
		}
	}

	go io.Copy(resp.Conn, os.Stdin)
	io.Copy(os.Stdout, resp.Reader)

	return nil
}
