package command

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yusufRahmatullah/game_save/repository"
)

// RootCommand handles arguments and execute corresponding
// sub-commands to interact with repositories
type RootCommand struct {
	rootCmd       *cobra.Command
	OSRepository  repository.IOSRepository
	GitRepository repository.IGitRepository
}

// NewRootCommand instantiate new RootCommand
// requires OSRepository and GitRepository
func NewRootCommand(osRepository repository.IOSRepository, gitRepository repository.IGitRepository) *RootCommand {
	root := RootCommand{
		rootCmd: &cobra.Command{
			Use:   "gamesave [global-flags] <command> [local-flags] <arguments>",
			Short: "Bring Game's save data to cloud",
			Long: `Synchronize game save data to cloud (git) by
					specifying game name as folder in git
					repository`,
		},
		OSRepository:  osRepository,
		GitRepository: gitRepository,
	}
	root.rootCmd.AddCommand(addCommand)
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
