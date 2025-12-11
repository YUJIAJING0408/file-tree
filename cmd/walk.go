package cmd

import (
	"encoding/json"
	"fmt"
	fileTree "github.com/YUJIAJING0408/file-tree"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
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
		tempDir, _ := os.MkdirTemp("", "file-tree-")
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {

			}
		}(tempDir)
		// web-ui for json
		addr := host + ":" + strconv.Itoa(port)
		if webFlag {
			go fileTree.NewTreeMapHttp(addr, tempDir)
		}

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
		// Ignore
		rules, ok := fileTree.ReadIgnore(ignorePath)
		if !ok {
			rules = make(fileTree.RuleList, 0)
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
		var countFile *fileTree.CountFile = nil
		if countFlag {
			countFile = fileTree.NewCountFile(6)
		}
		var done = make(chan struct{})
		go func(d chan struct{}) {
			err := rootDir.WalkSync(0, maxSyncDepth, countFile, rules)
			if err != nil {
				fmt.Println(err)
			}
			d <- struct{}{}
		}(done)
	loop:
		for {
			select {
			case <-done:
				break loop // break for
			default:
				time.Sleep(time.Millisecond * 100) //
				fmt.Printf("\rFileTree is walking '%s' .", dirPath)
			}
		}

		var treeBytes []byte
		switch outType {
		case "json":
			treeBytes, err = json.Marshal(&rootDir)
		case "yaml":
			treeBytes, err = yaml.Marshal(&rootDir)
		default:
			panic(fmt.Sprintf("unknown output type: %s", outType))
		}
		fileName := fmt.Sprintf("%s-%s", rootDir.Name, time.Now().Format("20060102150405"))
		outputFilePath := filepath.Join(output, fileName)
		fmt.Println(fmt.Sprintf("FileTree will be output to the %s directory in %s format.\nFullPath is %s", output, outType, outputFilePath))

		if countFlag {
			mix := countFile.Mix()
			go func() {
				// txt 输出
				//cf, err := os.Create(fmt.Sprintf("%s.txt", outputFilePath))
				//if err != nil {
				//	return
				//}
				//defer func(cf *os.File) {
				//	err := cf.Close()
				//	if err != nil {
				//		return
				//	}
				//}(cf)
				//var buf bytes.Buffer
				//for s, info := range mix {
				//	buf.WriteString(fmt.Sprintf("\t%s\t%d\t%s\n", s, info.Count, fileTree.ByteString(info.Size)))
				//}
				//_, err = cf.Write(buf.Bytes())
				// csv 输出
				var tmp = make([][]string, 0, len(mix))
				for s, info := range mix {
					tmp = append(tmp, []string{s, strconv.Itoa(int(info.Count)), fileTree.ByteString(info.Size)})
				}
				NewCSV(fmt.Sprintf("%s.csv", outputFilePath), []string{"suffix", "file_count", "file_size"}, tmp)
			}()
		}

		if webFlag {
			file, err := os.Create(fmt.Sprintf("%s.%s", filepath.Join(tempDir, fileName), "json"))
			if err != nil {
				return err
			}
			defer file.Close()
			if outType == "json" {
				file.Write(treeBytes)
			} else {
				tmp, _ := json.Marshal(&rootDir)
				file.Write(tmp)
			}
			fmt.Printf("\nFileTreeWeb: 'http://%s/treeMap?file_name=%s'.\n", addr, fileName)
		}

		saveFile, err := os.Create(fmt.Sprintf("%s.%s", outputFilePath, outType))
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(saveFile)
		saveFile.Write(treeBytes)
		if webFlag {
			for {
				time.Sleep(time.Millisecond * 100)
			}
		}
		return nil
	},
}

var (
	dirPath      string
	outType      string
	output       string
	ignorePath   string
	host         string
	port         int
	countFlag    bool
	webFlag      bool
	maxSyncDepth uint8
)

func init() {
	walkCmd.Flags().StringVarP(&dirPath, "path", "p", ".", "absolute path")
	_ = walkCmd.MarkFlagRequired("path")
	walkCmd.Flags().StringVarP(&outType, "type", "t", "json", "output type [json|yaml]")
	walkCmd.Flags().StringVarP(&output, "output", "o", ".", "output path | Default current directory ")
	walkCmd.Flags().StringVarP(&ignorePath, "ignorePath", "i", ".treeignore", "ignore path")
	walkCmd.Flags().Uint8VarP(&maxSyncDepth, "maxSyncDepth", "m", 8, "Asynchronous traversal will only be initiated when the folder depth is less than maxSyncDepth")
	walkCmd.Flags().BoolVarP(&countFlag, "count", "c", false, "count all things. ")
	walkCmd.Flags().BoolVarP(&webFlag, "web", "w", false, "show res by web")
	walkCmd.Flags().StringVarP(&host, "host", "H", "localhost", "host")
	walkCmd.Flags().IntVarP(&port, "port", "P", 8080, "port")
	rootCmd.AddCommand(walkCmd)
}

func isRoot(path string) bool {
	path = filepath.Clean(path)
	return filepath.Dir(path) == path
}
