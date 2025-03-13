package utils

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"os"
	"strings"
)

// ListFiles return all files in the path.
func ListFiles(path string) ([]os.DirEntry, error) {
	if !fileutil.IsExist(path) {
		return []os.DirEntry{}, nil
	}

	return os.ReadDir(path)
}

// FindFile 向上级联查找文件
// filename 文件名
// upCascadeSize 搜索的层数
func FindFile(filename string, upCascadeSize int) string {
	if upCascadeSize <= 0 {
		upCascadeSize = 1
	}
	for i := 0; i < upCascadeSize; i++ {
		if _, err := os.Stat(filename); err == nil {
			return filename
		} else {
			filename = "../" + filename
		}
	}
	return ""
}

func AbsolutePath(path ...string) string {
	var (
		pathList = make([]string, 0, len(path))
	)
	for _, v := range path {
		if v == "" {
			continue
		}
		pathList = append(pathList, strings.Trim(v, "/"))
	}
	return "/" + strings.Join(pathList, "/")
}
