package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Long:  "Print the version number of fileTree",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fileTree v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
