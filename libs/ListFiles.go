package libs

import (
	"os"
	"path/filepath"
	"strings"
)

func ListFiles(path string, ext []string, recursive bool) []string {
	var files []string
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			if !recursive {
				return filepath.SkipDir
			}
		} else {
			for _, e := range ext {
				if strings.HasSuffix(path, e) {
					files = append(files, path)
				}
			}
		}
		return nil
	})
	return files
}