/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configGenerateCmd represents the configGenerate command
var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate plugin files for shell scripts",
	Long:  `Generate plugin files for shell scripts`,
	Run: func(cmd *cobra.Command, args []string) {
		NewBashPluginConfigImpl().GeneratePluginConfig()
		NewZshPluginConfigImpl().GeneratePluginConfig()

		fmt.Println("Run source <(kubectl plugin_completion plugin_completion $SHELL_NAME) to update kubectl plugin completion")

	},
}

func init() {
	configCmd.AddCommand(configGenerateCmd)
}
