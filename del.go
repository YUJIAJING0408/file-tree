package fileTree

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

type DelFail struct {
	Path string `json:"path"`
	Msg  string `json:"msg"`
}

// PathDepth Calculate path depth through delimiter
func PathDepth(path string) int {
	return strings.Count(path, string(os.PathSeparator))
}

func DelSync(paths []string) (failed []DelFail) {
	// Traverse through it and sort the path according to its deeper level and higher ranking
	sort.Slice(paths, func(i, j int) bool {
		return PathDepth(paths[i]) < PathDepth(paths[j])
	})
	// Concurrency control
	concurrency := runtime.NumCPU() * 2
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, path := range paths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			err := os.RemoveAll(path) // 文件或目录都能删
			if err != nil {
				mu.Lock()
				failed = append(failed, DelFail{
					Path: path,
					Msg:  err.Error(),
				})
				mu.Unlock()
			} else {
				fmt.Println("deleted:", path)
			}
		}()
	}
	wg.Wait()

	if len(failed) > 0 {
		fmt.Println("\n—— 删除失败列表 ——")
		for _, e := range failed {
			fmt.Println(e)
		}
	}
	return failed
}
