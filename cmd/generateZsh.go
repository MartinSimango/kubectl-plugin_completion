/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateZshCmd represents the generateZsh command
var generateZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate completion script for zsh",
	Long:  `Generate completion script for zsh`,
	Run: func(cmd *cobra.Command, args []string) {
		script, _ := NewZshPluginConfigImpl().GenerateCompletionScript()
		fmt.Println(script)
	},
}

func init() {
	pluginCompletionCmd.AddCommand(generateZshCmd)

}
