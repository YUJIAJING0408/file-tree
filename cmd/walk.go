package cmd

import (
	"fmt"
	fileTree "github.com/YUJIAJING0408/file-tree"
	"github.com/spf13/cobra"
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
		var rd RootDir
		var tk *fileTree.TopK
		if topK > 0 {
			rd.topK = fileTree.NewTopK(int(topK))
		}
		if checkDuplicateFiles {
			rd.duplicateFile = fileTree.NewDuplicateFiles()
		}

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
		rd.ignoreRules, _ = fileTree.ReadIgnore(ignorePath)
		// If the output path does not exist, generate one
		if err = os.MkdirAll(output, 0777); err != nil {
			return err
		}
		name := stat.Name()
		if isRoot(dirPath) {
			name = dirPath[0:1]
		}
		rd.dir = fileTree.Dir{
			Name:     name,
			FullPath: dirPath,
			Type:     fileTree.TypeDir,
			Perm:     uint16(stat.Mode().Perm()),
		}
		fmt.Println(fmt.Sprintf("FileTree will walk '%s' .", dirPath))
		if countFlag {
			rd.countFile = fileTree.NewCountFile(6)
		}
		var done = make(chan struct{})
		go func(d chan struct{}) {
			err := rd.Walk()
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

		fileName := fmt.Sprintf("%s-%s", rd.dir.Name, time.Now().Format("20060102150405"))
		outputFilePath := filepath.Join(output, fileName)
		fmt.Println(fmt.Sprintf("FileTree will be output to the %s directory in %s format.\nFullPath is %s", output, outType, outputFilePath))

		if countFlag {
			mix := rd.countFile.Mix()
			go func() {
				// csv 输出
				var tmp = make([][]string, 0, len(mix))
				for s, info := range mix {
					tmp = append(tmp, []string{s, strconv.Itoa(int(info.Count)), fileTree.ByteString(info.Size)})
				}
				NewCSV(fmt.Sprintf("%s.csv", outputFilePath), []string{"suffix", "file_count", "file_size"}, tmp)
			}()
		}

		if webFlag {
			SaveToJson(fmt.Sprintf("%s.%s", filepath.Join(tempDir, fileName), "json"), &rd.dir)
			fmt.Printf("\nFileTreeWeb: 'http://%s/treeMap?file_name=%s'.\n", addr, fileName)
		}
		path := fmt.Sprintf("%s.%s", outputFilePath, outType)
		if outType == "json" {
			SaveToJson(path, &rd.dir)
		} else {
			SaveToYaml(path, &rd.dir)
		}

		if err != nil {
			return err
		}
		if webFlag {
			for {
				time.Sleep(time.Millisecond * 100)
			}
		}
		if topK > 0 && tk != nil {
			fileTree.FileHead(tk.TopKSorted()).Println()
		}

		if checkDuplicateFiles {
			ret := rd.duplicateFile.Check()
			fmt.Println(ret)
		}
		return nil
	},
}

var (
	dirPath             string
	outType             string
	output              string
	ignorePath          string
	host                string
	port                int
	countFlag           bool
	webFlag             bool
	maxSyncDepth        uint8
	topK                uint16
	checkDuplicateFiles bool
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
	walkCmd.Flags().Uint16VarP(&topK, "topK", "", 0, "topK size")
	walkCmd.Flags().BoolVarP(&checkDuplicateFiles, "checkDuplicateFiles", "C", false, "check duplicate files")
	rootCmd.AddCommand(walkCmd)
}

func isRoot(path string) bool {
	path = filepath.Clean(path)
	return filepath.Dir(path) == path
}

type RootDir struct {
	maxSyncDepth  uint8
	dir           fileTree.Dir
	countFile     *fileTree.CountFile
	ignoreRules   fileTree.RuleList
	topK          *fileTree.TopK
	duplicateFile *fileTree.DuplicateFiles
}

func (rd *RootDir) Walk() error {
	return rd.dir.WalkSync(0, rd.maxSyncDepth, rd.countFile, rd.ignoreRules, rd.topK, rd.duplicateFile)
}
