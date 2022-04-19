/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// cleanBashCmd represents the cleanBash command
var cleanBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Clear the bash config plugin file",
	Long:  `Clear the bash config plugin file`,
	Run: func(cmd *cobra.Command, args []string) {
		NewBashPluginConfigImpl().WriteCleanConfig()
	},
}

func init() {
	cleanCmd.AddCommand(cleanBashCmd)
}
