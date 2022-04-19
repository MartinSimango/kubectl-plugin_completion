/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// cleanAllCmd represents the cleanAll command
var cleanAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Clear plugin config file for all shells",
	Long:  `Clear plugin config file for all shells`,
	Run: func(cmd *cobra.Command, args []string) {
		NewBashPluginConfigImpl().WriteCleanConfig()
		NewZshPluginConfigImpl().WriteCleanConfig()
	},
}

func init() {
	cleanCmd.AddCommand(cleanAllCmd)
}
