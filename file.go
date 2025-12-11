package fileTree

import (
	"container/heap"
	"fmt"
	"sort"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

type File struct {
	Name     string `json:"name" yaml:"name"`
	FullPath string `json:"path" yaml:"full_path"`
	Size     int64  `json:"value" yaml:"size"`
	Type     uint8  `json:"type" yaml:"type"`
	Perm     uint16 `json:"perm" yaml:"perm"`
	Suffix   string `json:"suffix" yaml:"suffix"`
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetFullPath() string {
	return f.FullPath
}

func (f *File) GetSize() int64 {
	return f.Size
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) String() string {
	//TODO implement me
	panic("implement me")
}

func (f *File) Print(i int) {
	fmt.Println("*****************************************************")
	fmt.Printf("** Name: %s,\n** FullPath: %s,\n** Size: %s,\n** Type: %d,\n** Perm: %d,\n** Suffix: %s\n", f.Name, f.FullPath, ByteString(f.Size), f.Type, f.Perm, f.Suffix)
}

type FileHead []File

func (h FileHead) Len() int {
	return len(h)
}
func (h FileHead) Less(i, j int) bool {
	return h[i].Size < h[j].Size // 小顶堆
}

func (h FileHead) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *FileHead) Push(x interface{}) {
	*h = append(*h, x.(File))
}

func (h *FileHead) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h FileHead) Println() {
	for _, file := range h {
		file.Print(0)
	}
}

type TopK struct {
	k int
	h *FileHead
}

func (t *TopK) Enabled() bool {
	return t.k > 0
}

func NewTopK(k int) *TopK {
	h := &FileHead{}
	heap.Init(h)
	return &TopK{
		k: k,
		h: h,
	}
}

func (t *TopK) Push(f File) {
	if t.h.Len() < t.k {
		heap.Push(t.h, f)
		return
	}

	if f.Size > (*t.h)[0].Size {
		heap.Pop(t.h)
		heap.Push(t.h, f)
	}
}

func (t *TopK) TopK() []File {
	return *t.h
}

func (t *TopK) TopKSorted() []File {
	res := make([]File, t.h.Len())
	copy(res, *t.h)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Size > res[j].Size
	})
	return res
}
