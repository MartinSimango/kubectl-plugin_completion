/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify plugin config files for specified shell (bash or zsh)",
	Long: `Modify plugin config files for specified shell (bash or zsh) using subcommands like "kubectl plugin_completion config clean zsh"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
