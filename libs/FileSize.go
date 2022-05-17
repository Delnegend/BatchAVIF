package libs

import (
	"os"
	"fmt"
)

func FileSize(path string) (int64) {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}
func HumanReadableSize(size int64) string {
	size_str := ""
	if size > 1024*1024*1024 {
		size_str = fmt.Sprintf("%d GB", size/1024/1024/1024)
	} else if size > 1024*1024 {
		size_str = fmt.Sprintf("%d MB", size/1024/1024)
	} else if size > 1024 {
		size_str = fmt.Sprintf("%d KB", size/1024)
	} else {
		size_str = fmt.Sprintf("%d B", size)
	}
	return size_str
}