package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

var rootCmd = &cobra.Command{
	Use:   "fileTree",
	Short: "fileTree can quickly generate file trees",
	Long:  `fileTree can quickly generate file trees in various formats as needed`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run fileTree ...")
		if len(args) == 0 {
			fmt.Println("what about try 'fileTree help'?")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
