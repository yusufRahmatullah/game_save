package command

import (
	"fmt"
	"io"
	"os"

	"github.com/yusufRahmatullah/game_save/service"

	"github.com/spf13/cobra"
)

var (
	appVersion  string
	rootService service.IService
)

// RootCommand handles arguments and execute corresponding
// sub-commands to interact with service
type RootCommand struct {
	rootCmd *cobra.Command
}

// NewRootCommand instantiate new RootCommand
// requires OSRepository and GitRepository
func NewRootCommand(serv service.IService) *RootCommand {
	root := RootCommand{
		rootCmd: &cobra.Command{
			Use:   "gamesave [global-flags] <command> [local-flags] <arguments>",
			Short: "Bring Game's save data to cloud",
			Long: `Synchronize game save data to cloud (git) by
					specifying game name as folder in git
					repository`,
		},
	}
	rootService = serv
	root.rootCmd.AddCommand(addCommand)
	root.rootCmd.AddCommand(initCommand)
	root.rootCmd.AddCommand(loadCommand)
	root.rootCmd.AddCommand(saveCommand)
	root.rootCmd.AddCommand(setPathCommand)
	root.rootCmd.AddCommand(versionCommand)
	return &root
}

// Parse receive arguments as list of string and
// write the result to output
func (c *RootCommand) Parse(args []string, output io.Writer) error {
	c.rootCmd.SetArgs(args)
	c.rootCmd.SetOutput(output)
	return c.rootCmd.Execute()
}

// Run executes the main app to be callable
// from command line
func (c *RootCommand) Run() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// SetVersion set RootComand verison
func (c *RootCommand) SetVersion(version string) {
	appVersion = version
}
