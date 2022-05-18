package libs

import (
	"fmt"
	"os"
)

func FileSize(path string) float64 {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return 0
	}
	return float64(fi.Size())
}

func HumanReadableSize(size float64) string {
	size_str := ""
	if size > 1024*1024*1024 {
		size_str = fmt.Sprintf("%.2f GB", size/1024/1024/1024)
	} else if size > 1024*1024 {
		size_str = fmt.Sprintf("%.2f MB", size/1024/1024)
	} else if size > 1024 {
		size_str = fmt.Sprintf("%.2f KB", size/1024)
	} else {
		size_str = fmt.Sprintf("%.f B", size)
	}
	return size_str
}

func ReportFileSize(old float64, new float64) {
	fmt.Printf(
		"Original: %s | Converted: %s ~ %.2f%%\n",
		HumanReadableSize(old),
		HumanReadableSize(new),
		float64(new)/float64(old)*100,
	)
}
