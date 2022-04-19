/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// pluginCompletionCmd represents the generateCompletion command
var pluginCompletionCmd = &cobra.Command{
	Use:     "plugin-completion",
	Short:   "Generates completion script for plugins for a specified shell (bash or zsh)",
	Long:    `Generates completion script for plugins for a specified shell (bash or zsh)`,
	Example: "  kubectl plugin_completion plugin-completion zsh",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(pluginCompletionCmd)
}
