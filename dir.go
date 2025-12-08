package fileTree

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

var CountFileNum = sync.Map{}
var CountFileSize = sync.Map{}

type Dir struct {
	sync.Mutex `json:"-" yaml:"-"`
	Name       string `json:"name" yaml:"name"`
	FullPath   string `json:"full_path" yaml:"full_path"`
	Size       int64  `json:"size" yaml:"size"`
	Type       uint16 `json:"type" yaml:"type"`
	Perm       uint16 `json:"perm" yaml:"perm"`
	Children   []any  `json:"children" yaml:"children"`
}

func (d *Dir) Walk() error {
	items, err := os.ReadDir(d.FullPath)
	if err != nil {
		return err
	}
	// 遍历
	for _, item := range items {
		info, err := item.Info()
		if err != nil {
			return err
		}
		if info.IsDir() {
			// 构建Dir
			dir := &Dir{
				Name:     info.Name(),
				FullPath: filepath.Join(d.FullPath, info.Name()),
				Type:     TypeDir,
				Perm:     uint16(info.Mode().Perm()),
			}
			err = dir.Walk()
			if err != nil {
				return err
			}
			d.Lock()
			d.Size += dir.Size
			d.Children = append(d.Children, dir)
			d.Unlock()
		} else {
			split := strings.Split(info.Name(), ".")
			suffix := split[len(split)-1]
			// file or lnk
			if suffix == "lnk" {
				link := &Link{
					Name:     info.Name(),
					FullPath: filepath.Join(d.FullPath, info.Name()),
					Type:     TypeLink,
					Perm:     uint16(info.Mode().Perm()),
					Size:     info.Size(),
				}
				link.LinkTo, err = getLnkTargetPath(link.FullPath)
				if err != nil {
					return err
				}
				d.Lock()
				d.Children = append(d.Children, link)
				d.Size += link.Size
				d.Unlock()
			} else {
				file := &File{
					Name:     info.Name(),
					FullPath: filepath.Join(d.FullPath, info.Name()),
					Type:     TypeFile,
					Perm:     uint16(info.Mode().Perm()),
					Suffix:   suffix,
					Size:     info.Size(),
				}
				d.Lock()
				d.Children = append(d.Children, file)
				d.Size += file.Size
				d.Unlock()

			}
		}
	}
	return nil
}

func (d *Dir) WalkSync(depth uint8, syncMaxDepth uint8) error {
	if depth < syncMaxDepth {
		// 遍历FullPath
		items, err := os.ReadDir(d.FullPath)
		if err != nil {
			return err
		}
		// 使用协程
		wg := &sync.WaitGroup{}
		// 遍历
		for _, item := range items {
			info, err := item.Info()
			if err != nil {
				return err
			}
			if info.IsDir() {
				// 构建Dir
				dir := &Dir{
					Name:     info.Name(),
					FullPath: filepath.Join(d.FullPath, info.Name()),
					Type:     TypeDir,
					Perm:     uint16(info.Mode().Perm()),
				}
				wg.Go(func() {
					err := dir.WalkSync(depth+1, syncMaxDepth)
					if err != nil {
						return
					}
					d.Lock()
					d.Size += dir.Size
					d.Children = append(d.Children, dir)
					d.Unlock()
				})
			} else {
				split := strings.Split(info.Name(), ".")
				suffix := split[len(split)-1]
				// file or lnk
				if suffix == "lnk" {
					link := &Link{
						Name:     info.Name(),
						FullPath: filepath.Join(d.FullPath, info.Name()),
						Type:     TypeLink,
						Perm:     uint16(info.Mode().Perm()),
						Size:     info.Size(),
					}
					link.LinkTo, err = getLnkTargetPath(link.FullPath)
					if err != nil {
						return err
					}
					d.Lock()
					d.Children = append(d.Children, link)
					d.Size += link.Size
					d.Unlock()
				} else {
					file := &File{
						Name:     info.Name(),
						FullPath: filepath.Join(d.FullPath, info.Name()),
						Type:     TypeFile,
						Perm:     uint16(info.Mode().Perm()),
						Suffix:   suffix,
						Size:     info.Size(),
					}
					d.Lock()
					d.Children = append(d.Children, file)
					d.Size += file.Size
					d.Unlock()
				}
			}
		}
		wg.Wait()
		return nil
	} else {
		// 层级已经很深了
		err := d.Walk()
		if err != nil {
			return err
		}
		return nil
	}
}
