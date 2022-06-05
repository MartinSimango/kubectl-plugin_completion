/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// listZshCmd represents the listZsh command
var listZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "List all plugins in zsh config file",
	Long:  `List all plugins in zsh config file`,
	Run: func(cmd *cobra.Command, args []string) {
		NewZshPluginConfigImpl().PrintAllPlugins()

	},
}

func init() {
	listCmd.AddCommand(listZshCmd)

}
