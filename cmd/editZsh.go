/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// editZshCmd represents the editZsh command
var editZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Edit plugin in zsh config file",
	Long:  `Edit plugin in zsh config file`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := cmd.Flag("plugin").Value.String()
		description := cmd.Flag("description").Value.String()
		completionFunctionName := cmd.Flag("completionFunction").Value.String()

		descriptionSet := cmd.Flag("description").Changed
		completionFunctionSet := cmd.Flag("completionFunction").Changed

		err := NewZshPluginConfigImpl().EditPlugin(pluginName, completionFunctionName, description, completionFunctionSet, descriptionSet)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("zsh config file edited")
	},
}

func init() {
	editCmd.AddCommand(editZshCmd)

	editZshCmd.Flags().StringP("plugin", "p", "", "name of the plugin")
	editZshCmd.Flags().StringP("description", "d", "", "description of the plugin")
	editZshCmd.Flags().StringP("completion-function", "f", "", "completion function name of the plugin")
	editZshCmd.MarkFlagRequired("plugin")
}
