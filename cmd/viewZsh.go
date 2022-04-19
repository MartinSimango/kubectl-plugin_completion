/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// viewZshCmd represents the viewZsh command
var viewZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Display the contents of the zsh plugin config file",
	Long:  `Display the contents of the zsh plugin config file`,
	Run: func(cmd *cobra.Command, args []string) {
		NewZshPluginConfigImpl().DisplayConfig()
	},
}

func init() {
	viewCmd.AddCommand(viewZshCmd)
}
