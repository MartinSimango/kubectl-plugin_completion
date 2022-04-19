/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display the contents of all plugin config files",
	Long: `Displays the contents of all plugin config files or just the plugin config file for specified shell (bash or zsh).

For example to view the plugin config of the zsh shell, enter the following command:
"kubectl plugin_completion config view zsh"`,
	Run: func(cmd *cobra.Command, args []string) {
		displayConfig()
	},
}

func init() {
	configCmd.AddCommand(viewCmd)

}

func displayConfig() {
	NewBashPluginConfigImpl().DisplayConfig()
	fmt.Println("\n---\n")
	NewZshPluginConfigImpl().DisplayConfig()
}
