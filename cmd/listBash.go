/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// listZshCmd represents the listZsh command
var listBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "List all plugins in bash config file",
	Long:  `List all plugins in bash config file`,
	Run: func(cmd *cobra.Command, args []string) {
		NewBashPluginConfigImpl().PrintAllPlugins()

	},
}

func init() {
	listCmd.AddCommand(listBashCmd)

}
