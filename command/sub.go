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
	RunE: func(cmd *cobra.Command, args []string) error {
		gameName := args[0]
		err := rootService.AddConfig("game_name", gameName)
		if err != nil {
			return err
		}
		err = rootService.PrepareGame()
		if err != nil {
			return err
		}
		fmt.Printf("Your game is %v\n", gameName)
		return nil
	},
}

var initCommand = &cobra.Command{
	Use:   "init <git repo URL>",
	Short: "Initialize GameSave in this machine",
	Long: `Initialize GameSave in this machine by clone
			given Git Repository URL. Ensure the repository
			is exist`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL := args[0]
		return rootService.InitGitRepo(repoURL)
	},
}

var loadCommand = &cobra.Command{
	Use:   "load",
	Short: "Load game",
	Long:  `Load game by synchronize save from the cloud`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := rootService.PrepareGame(); err != nil {
			return err
		}
		return rootService.LoadGame()
	},
}

var saveCommand = &cobra.Command{
	Use:   "save",
	Short: "Save game",
	Long:  "Save game by synchronize save to the cloud",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootService.SaveGame()
	},
}

var setPathCommand = &cobra.Command{
	Use:   "set-path <game save path>",
	Short: "Set game save path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		savePath := args[0]
		return rootService.AddConfig("save_path", savePath)
	},
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Show gamesave version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()
		out.Write([]byte(appVersion))
		fmt.Println("")
	},
}
