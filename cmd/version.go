package cmd

import (
	"fmt"
	fileTree "github.com/YUJIAJING0408/file-tree"
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
		fmt.Println("fileTree v" + fileTree.VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
