package fileTree

import (
	"bufio"
	"os"
	"strings"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

type Rule interface {
	// Ignore true mean no pass
	Ignore(node Node) bool
}

type SuffixRule struct {
	suffix string
}

func (s SuffixRule) Ignore(node Node) bool {
	if node.IsDir() {
		// pass this rule
		return false
	}
	return strings.HasSuffix(node.GetName(), s.suffix)
}

type DirectoryRule struct {
	path string
}

func (d DirectoryRule) Ignore(node Node) bool {
	return strings.HasPrefix(node.GetFullPath(), d.path)
}

type FilePathRule struct {
	path string
}

func (f FilePathRule) Ignore(node Node) bool {
	return node.GetFullPath() == f.path
}

type FileNameRule struct {
	name string
}

func (f FileNameRule) Ignore(node Node) bool {
	return node.GetName() == f.name
}

type FileSizeRule struct {
	size        int64
	greaterThan bool
}

func (f FileSizeRule) Ignore(node Node) bool {
	if f.greaterThan {
		return node.GetSize() > f.size
	} else {
		return node.GetSize() <= f.size
	}
}

type RuleList []Rule

func (r RuleList) Ignore(node Node) (res bool) {
	for _, rule := range r {
		if rule.Ignore(node) {
			return true
		}
	}
	return false
}

func ReadIgnore(path string) (rules RuleList, ok bool) {
	if strings.HasSuffix(path, ".treeignore") {
		file, err := os.Open(path)
		if err != nil {
			return nil, false
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				// pass empty
				continue
			}
			if strings.HasPrefix(line, ">") {
				res, err := StringToByte(line[1:])
				if err != nil {
					return nil, false
				}
				rules = append(rules, FileSizeRule{size: int64(res), greaterThan: true})
				continue
			} else if strings.HasPrefix(line, "<") {
				res, err := StringToByte(line[1:])
				if err != nil {
					return nil, false
				}
				rules = append(rules, FileSizeRule{size: int64(res), greaterThan: false})
				continue
			}

			if strings.HasPrefix(line, "*.") {
				rules = append(rules, SuffixRule{suffix: line[2:]})
				continue
			}
			info, err := os.Stat(line)
			if err != nil {
				rules = append(rules, FileNameRule{name: line})
			} else {
				if info.IsDir() {
					rules = append(rules, DirectoryRule{path: line})
				} else {
					rules = append(rules, FilePathRule{path: line})
				}
			}

		}
	} else {
		return nil, false
	}
	return rules, true
}
