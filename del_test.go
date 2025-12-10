package fileTree

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path"
	"testing"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

func MakeTestDirAndFile() (paths []string) {
	dir1 := "test1"
	dir2 := "test2"
	paths = append(paths, dir1, path.Join(dir1, "file1.txt"), path.Join(dir1, "file2.txt"), path.Join(dir1, dir2), path.Join(dir1, dir2, "file3.txt"), path.Join(dir1, dir2, "file4.txt"))
	os.Mkdir(dir1, os.ModePerm)
	file1, _ := os.Create(path.Join(dir1, "file1.txt"))
	file1.Close()
	file2, _ := os.Create(path.Join(dir1, "file2.txt"))
	file2.Close()
	os.Mkdir(path.Join(dir1, dir2), os.ModePerm)
	file3, _ := os.Create(path.Join(dir1, dir2, "file3.txt"))
	file3.Close()
	file4, _ := os.Create(path.Join(dir1, dir2, "file4.txt"))
	file4.Close()
	rand.Shuffle(len(paths), func(i, j int) { paths[i], paths[j] = paths[j], paths[i] })
	return paths
}

func TestDel(t *testing.T) {
	paths := MakeTestDirAndFile()
	fmt.Println("paths:", paths)
	fmt.Println(DelSync(paths))
}
