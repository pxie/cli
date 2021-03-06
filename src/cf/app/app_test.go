package app

import (
	"cf/api"
	"cf/commands"
	"cf/configuration"
	"cf/net"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
	"strings"
	testassert "testhelpers/assert"
	testconfig "testhelpers/configuration"
	testmanifest "testhelpers/manifest"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func availableCmdNames() (names []string) {
	reqFactory := &testreq.FakeReqFactory{}
	cmdRunner := commands.NewRunner(nil, reqFactory)
	app, _ := NewApp(cmdRunner)

	for _, cliCmd := range app.Commands {
		if cliCmd.Name != "help" {
			names = append(names, cliCmd.Name)
		}
	}
	return
}

type FakeRunner struct {
	cmdFactory commands.Factory
	t          *testing.T
	cmdName    string
}

func (runner *FakeRunner) RunCmdByName(cmdName string, c *cli.Context) (err error) {
	_, err = runner.cmdFactory.GetByCmdName(cmdName)
	if err != nil {
		runner.t.Fatal("Error instantiating command with name", cmdName)
		return
	}
	runner.cmdName = cmdName
	return
}

func TestCommands(t *testing.T) {
	for _, cmdName := range availableCmdNames() {
		ui := &testterm.FakeUI{}
		config := &configuration.Configuration{}
		configRepo := testconfig.FakeConfigRepository{}
		manifestRepo := &testmanifest.FakeManifestRepository{}

		repoLocator := api.NewRepositoryLocator(config, configRepo, map[string]net.Gateway{
			"auth":             net.NewUAAGateway(),
			"cloud-controller": net.NewCloudControllerGateway(),
			"uaa":              net.NewUAAGateway(),
		})

		cmdFactory := commands.NewFactory(ui, config, configRepo, manifestRepo, repoLocator)
		cmdRunner := &FakeRunner{cmdFactory: cmdFactory, t: t}
		app, _ := NewApp(cmdRunner)
		app.Run([]string{"", cmdName})

		assert.Equal(t, cmdRunner.cmdName, cmdName)
	}
}

func TestUsageIncludesCommandName(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{}
	cmdRunner := commands.NewRunner(nil, reqFactory)
	app, _ := NewApp(cmdRunner)

	for _, cmd := range app.Commands {
		assert.Contains(t, strings.Split(cmd.Usage, "\n")[0], cmd.Name)
	}
}

func TestPushCommandHelpOutput(t *testing.T) {
	reqFactory := &testreq.FakeReqFactory{}
	cmdRunner := commands.NewRunner(nil, reqFactory)
	app, _ := NewApp(cmdRunner)

	var updateCommand, pushCommand cli.Command
	for _, cmd := range app.Commands {
		if cmd.Name == "push" {
			pushCommand = cmd
		} else if cmd.Name == "update-buildpack" {
			updateCommand = cmd
		}
	}

	flags := []string{}
	for _, flag := range pushCommand.Flags {
		flags = append(flags, flag.String())
	}

	testassert.SliceContains(t, flags, testassert.Lines{
		{"-b \tCustom buildpack URL (e.g. https://github.com/heroku/heroku-buildpack-play.git)"},
	})

	flags = []string{}
	for _, flag := range updateCommand.Flags {
		flags = append(flags, flag.String())
	}

	testassert.SliceContains(t, flags, testassert.Lines{
		{"-i \tBuildpack position among other buildpacks"},
	})
}
