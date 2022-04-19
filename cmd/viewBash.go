/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// viewBashCmd represents the viewBash command
var viewBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Display the contents of the bash plugin config file",
	Long:  `Display the contents of the bash plugin config file`,
	Run: func(cmd *cobra.Command, args []string) {
		NewBashPluginConfigImpl().DisplayConfig()
	},
}

func init() {
	viewCmd.AddCommand(viewBashCmd)

}
