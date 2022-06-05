/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// editBashCmd represents the editBash command
var editBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Edit plugin in bash config file",
	Long:  `Edit plugin in bash config file`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := cmd.Flag("plugin").Value.String()
		description := cmd.Flag("description").Value.String()
		completionFunctionName := cmd.Flag("completionFunction").Value.String()

		descriptionSet := cmd.Flag("description").Changed
		completionFunctionSet := cmd.Flag("completionFunction").Changed
		err := NewBashPluginConfigImpl().EditPlugin(pluginName, completionFunctionName, description, completionFunctionSet, descriptionSet)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("bash config file edited")
		fmt.Println("Please run 'source <(kubectl plugin_completion plugin-completion bash)' to ensure your changes take effect immediately.")

	},
}

func init() {
	editCmd.AddCommand(editBashCmd)

	editBashCmd.Flags().StringP("plugin", "g", "", "name of the plugin")
	editBashCmd.Flags().StringP("description", "d", "", "description of the plugin")
	editBashCmd.Flags().StringP("completion-function", "f", "", "completion function name of the plugin")
	editBashCmd.MarkFlagRequired("plugin")

}
