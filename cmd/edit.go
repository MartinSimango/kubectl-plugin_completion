/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit plugin config for a specified shell (bash or zsh)",
	Long: `Edit plugin config for a specified shell (bash or zsh) using subcommands like:
"kubectl plugin_completion config edit zsh"`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	configCmd.AddCommand(editCmd)

}
