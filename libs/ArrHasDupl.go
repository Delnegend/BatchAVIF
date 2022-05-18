package libs

import (
	// "fmt"
	"path/filepath"
	"strings"
)

func ArrHasDupl(array_ []string) []string {
	array := make([]string, len(array_))
	copy(array, array_)
	var dupls []string
	for i := 0; i < len(array); i++ {
		for j := 0; j < len(array); j++ {
			if (rmExt(array[i]) == rmExt(array[j])) && (i != j) && !hasElem(dupls, array[j]) {
				dupls = append(dupls, array[j])
			}
		}
	}
	if len(dupls) > 0 {
		return dupls
	}
	return nil
}

func rmExt(name string) string {
	return strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
}

func hasElem(array_ []string, elem string) bool {
	array := make([]string, len(array_))
	copy(array, array_)
	for _, i := range array {
		if i == elem {
			return true
		}
	}
	return false
}