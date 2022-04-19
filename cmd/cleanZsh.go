/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// cleanZshCmd represents the cleanZsh command
var cleanZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Clear the zsh config plugin file",
	Long:  `Clear the zsh config plugin file`,
	Run: func(cmd *cobra.Command, args []string) {
		NewZshPluginConfigImpl().WriteCleanConfig()
	},
}

func init() {
	cleanCmd.AddCommand(cleanZshCmd)

}
