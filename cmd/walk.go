package cmd

import (
	"encoding/json"
	"fmt"
	fileTree "github.com/YUJIAJING0408/file-tree"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Traverse",
	Long:  "Traverse the input directory and its subdirectories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dirPath == "." {
			dirPath, _ = os.Getwd()
		}
		if output == "." {
			output, _ = os.Getwd()
		}
		// check dirPath and output
		stat, err := os.Stat(dirPath)
		if err != nil {
			return err
		} else {
			if !stat.IsDir() {
				return fmt.Errorf("%s is not a directory", dirPath)
			}
		}
		// If the output path does not exist, generate one
		if err = os.MkdirAll(output, 0777); err != nil {
			return err
		}
		name := stat.Name()
		if isRoot(dirPath) {
			name = dirPath[0:1]
		}
		var rootDir = fileTree.Dir{
			Name:     name,
			FullPath: dirPath,
			Type:     fileTree.TypeDir,
			Perm:     uint16(stat.Mode().Perm()),
		}
		fmt.Println(fmt.Sprintf("FileTree will walk '%s' .", dirPath))
		var done = make(chan struct{})
		go func(d chan struct{}) {
			err := rootDir.WalkSync(0, maxSyncDepth)
			if err != nil {
				fmt.Println(err)
			}
			d <- struct{}{}
		}(done)
		var count = 0
		var countSlice = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	loop:
		for {
			select {
			case <-done:
				break loop // break for
			default:
				time.Sleep(time.Millisecond * 100) //
				fmt.Printf("\rFileTree is walking '%s' [%s] .", dirPath, countSlice[count%len(countSlice)])
			}
		}
		var bytes []byte
		switch outType {
		case "json":
			bytes, err = json.Marshal(&rootDir)
		case "yaml":
			bytes, err = yaml.Marshal(&rootDir)
		default:
			panic(fmt.Sprintf("unknown output type: %s", outType))
		}
		outputFilePath := filepath.Join(output, fmt.Sprintf("[%s-%s].%s", rootDir.Name, time.Now().Format("20060102150405"), outType))
		fmt.Println(fmt.Sprintf("FileTree will be output to the %s directory in %s format.\nFullPath is %s", output, outType, outputFilePath))
		file, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)
		file.Write(bytes)
		return nil
	},
}

var (
	dirPath      string
	outType      string
	output       string
	count        bool
	maxSyncDepth uint8
)

func init() {
	walkCmd.Flags().StringVarP(&dirPath, "path", "p", ".", "absolute path")
	_ = walkCmd.MarkFlagRequired("path")
	walkCmd.Flags().StringVarP(&outType, "type", "t", "json", "output type [json|yaml]")
	walkCmd.Flags().StringVarP(&output, "output", "o", ".", "output path | Default current directory ")
	walkCmd.Flags().Uint8VarP(&maxSyncDepth, "maxSyncDepth", "m", 8, "Asynchronous traversal will only be initiated when the folder depth is less than maxSyncDepth")
	walkCmd.Flags().BoolVarP(&count, "count", "c", false, "count all things. ")

	rootCmd.AddCommand(walkCmd)
}

func isRoot(path string) bool {
	path = filepath.Clean(path)
	return filepath.Dir(path) == path
}
