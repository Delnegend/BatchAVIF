package libs

import (
	"strings"
	"path/filepath"
)

func ArrHasDupl(array_ []string) []string {
	array := make([]string, len(array_))
	copy(array, array_)
	var dupls []string
	for i := 0; i < len(array); i++ {
		for j := i; j < len(array); j++ {
			if rmExt(array[i]) == rmExt(array[j]) && (i != j) {
				if !InArr(dupls, array[i]) {
					dupls = append(dupls, array[i])
				}
				if !InArr(dupls, array[j]) {
					dupls = append(dupls, array[j])
				}
			}
		}
	}
	if len(dupls) > 0 {
		return dupls
	}
	return nil
}

func rmExt(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}