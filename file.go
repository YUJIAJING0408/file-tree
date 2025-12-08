package fileTree

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

type File struct {
	Name     string `json:"name" yaml:"name"`
	FullPath string `json:"full_path" yaml:"full_path"`
	Size     int64  `json:"size" yaml:"size"`
	Type     uint8  `json:"type" yaml:"type"`
	Perm     uint16 `json:"perm" yaml:"perm"`
	Suffix   string `json:"suffix" yaml:"suffix"`
}

func (f *File) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (f *File) GetFullPath() string {
	//TODO implement me
	panic("implement me")
}

func (f *File) GetSize() int64 {
	//TODO implement me
	panic("implement me")
}

func (f *File) IsDir() bool {
	//TODO implement me
	panic("implement me")
}

func (f *File) String() string {
	//TODO implement me
	panic("implement me")
}

func (f *File) Print(i int) {
	//TODO implement me
	panic("implement me")
}
