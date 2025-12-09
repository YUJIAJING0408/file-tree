package fileTree

import (
	"crypto/md5"
	"hash/maphash"
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

var seed maphash.Seed

func init() {
	seed = maphash.MakeSeed()
}

type FileInfo struct {
	Size  int64
	Count int64
}

func (f *FileInfo) Add(s int64) {
	f.Size += s
	f.Count++
}

type CountFile struct {
	shardBits uint32
	shardNum  uint32
	shardMask uint32
	data      []*ShardedMap
}

type ShardedMap struct {
	sync.Mutex
	data map[string]*FileInfo
}

func NewShardedMap() *ShardedMap {
	return &ShardedMap{
		data: make(map[string]*FileInfo),
	}
}

func shardIndex(key [16]byte, shardMask uint32) uint32 {
	var h maphash.Hash
	h.SetSeed(seed)
	_, _ = h.Write(key[:])
	return uint32(h.Sum64()) & shardMask
}

func NewCountFile(shardBit uint32) *CountFile {
	var shardNum uint32 = 1 << shardBit
	var d = make([]*ShardedMap, shardNum)
	for i := uint32(0); i < shardNum; i++ {
		d[i] = NewShardedMap()
	}
	return &CountFile{shardBit, shardNum, shardNum - 1, d}
}

func (cf *CountFile) Add(key string, size int64) {
	// hash
	hash := md5.Sum([]byte(key))
	index := shardIndex(hash, cf.shardMask)
	m := cf.data[index]
	m.Lock()
	file := m.data[key]
	if file == nil {
		// First Time
		file = &FileInfo{}
		m.data[key] = file
	}
	file.Add(size)
	m.Unlock()
}

func (cf *CountFile) Mix() map[string]*FileInfo {
	var m = make(map[string]*FileInfo)
	for _, shardMap := range cf.data {
		tmp := (*shardMap).data
		for key, value := range tmp {
			m[key] = value
		}
	}
	return m
}

// Dir Directory
type Dir struct {
	sync.Mutex `json:"-" yaml:"-"`
	Name       string `json:"name" yaml:"name"`
	FullPath   string `json:"full_path" yaml:"full_path"`
	Size       int64  `json:"size" yaml:"size"`
	Type       uint16 `json:"type" yaml:"type"`
	Perm       uint16 `json:"perm" yaml:"perm"`
	Children   []any  `json:"children" yaml:"children"`
}

// Walk The most normal depth first traversal
func (d *Dir) Walk(countFile *CountFile) error {
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
			err = dir.Walk(countFile)
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
			if countFile != nil {
				countFile.Add(suffix, info.Size())
			}
		}
	}
	return nil
}

func (d *Dir) WalkSync(depth uint8, syncMaxDepth uint8, countFile *CountFile) error {
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
					err := dir.WalkSync(depth+1, syncMaxDepth, countFile)
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
				if countFile != nil {
					countFile.Add(suffix, info.Size())
				}
			}
		}
		wg.Wait()
		return nil
	} else {
		// 层级已经很深了
		err := d.Walk(countFile)
		if err != nil {
			return err
		}
		return nil
	}
}

func (d *Dir) Find() {

}
