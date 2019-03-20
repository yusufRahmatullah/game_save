package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add <game name>",
	Short: "Add game name",
	Long:  "Specify a unique name of the current game",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Your game is %v\n", args[0])
	},
}
