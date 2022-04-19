/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear config files for specified shell (bash or zsh) or for all shells (by specifying all instead of a shell name)",
	Long: `Clear config files for specified shell (bash or zsh) or for all shells (by specifying all instead of a shell name).

	To clean the config file for the zsh terminal, you can enter the following command: 
	"kubectl plugin_completion clean zsh"`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	configCmd.AddCommand(cleanCmd)

}
