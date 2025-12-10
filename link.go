package fileTree

import (
	lnk "github.com/parsiya/golnk"
	"golang.org/x/text/encoding/simplifiedchinese"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

func getLnkTargetPath(lnkPath string) (string, error) {
	// 解析快捷方式文件
	link, err := lnk.File(lnkPath)
	if err != nil {
		return "", err
	}
	// 获取目标路径
	targetPath := link.LinkInfo.LocalBasePath
	if targetPath == "" {
		targetPath = link.LinkInfo.CommonPathSuffix
	}
	// 处理中文字符编码
	decoder := simplifiedchinese.GB18030.NewDecoder()
	targetPath, _ = decoder.String(targetPath)
	return targetPath, nil
}

type Link struct {
	Name     string `json:"name" yaml:"name"`
	FullPath string `json:"full_path" yaml:"full_path"`
	Size     int64  `json:"size" yaml:"size"`
	Type     uint8  `json:"type" yaml:"type"`
	Perm     uint16 `json:"perm" yaml:"perm"`
	LinkTo   string `json:"link_to" yaml:"link_to"`
}

func (l Link) GetName() string {
	return l.Name
}

func (l Link) GetFullPath() string {
	return l.FullPath
}

func (l Link) GetSize() int64 {
	return l.Size
}

func (l Link) IsDir() bool {
	return false
}

func (l Link) String() string {
	//TODO implement me
	panic("implement me")
}

func (l Link) Print(i int) {
	//TODO implement me
	panic("implement me")
}
