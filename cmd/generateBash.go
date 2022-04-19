/*
Copyright © 2022 Martin Simango <shukomango@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateBashCmd represents the generateBash command
var generateBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate completion script for bash",
	Long:  `Generate completion script for bash`,
	Run: func(cmd *cobra.Command, args []string) {
		script, _ := NewBashPluginConfigImpl().GenerateCompletionScript()
		fmt.Println(script)
	},
}

func init() {
	pluginCompletionCmd.AddCommand(generateBashCmd)
}
