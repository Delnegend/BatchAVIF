package libs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ProcessPreset(
	log *os.File,
	file string,
	ext_ []string,
	enc_ []string,
	fallback_ []string,
	repack_ []string,
) ([]string, []string, []string, []string, error) {
	var cf Config
	cf.ParseConfig()

	without_ext := file
	if !cf.Config.KeepOriginalExtension {
		without_ext = strings.TrimSuffix(file, filepath.Ext(file))
	}

	ext := make([]string, len(ext_))
	copy(ext, ext_)
	enc := make([]string, len(enc_))
	copy(enc, enc_)
	fallback := make([]string, len(fallback_))
	copy(fallback, fallback_)
	repack := make([]string, len(repack_))
	copy(repack, repack_)

	for i, p := range ext {
		if strings.Contains(p, "{{ input }}") {
			ext[i] = file
		}
		if strings.Contains(p, "{{ output }}") {
			ext[i] = file + ".y4m"
		}
	}
	for i, p := range enc {
		if strings.Contains(p, "{{ input }}") {
			enc[i] = strings.Replace(p, "{{ input }}", file+".y4m", -1)
		}
		if strings.Contains(p, "{{ output }}") {
			enc[i] = strings.Replace(p, "{{ output }}", file+".ivf", -1)
		}
		if strings.Contains(p, "{{ threads }}") {
			enc[i] = strings.Replace(p, "{{ threads }}", fmt.Sprintf("%d", MaxCPU()), -1)
		}
		if strings.Contains(p, "{{ width }}") || strings.Contains(p, "{{ height }}") {
			width, height, err := Dimension(log, file)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			if strings.Contains(p, "{{ width }}") {
				enc[i] = strings.Replace(p, "{{ width }}", fmt.Sprintf("%d", width), -1)
			}
			if strings.Contains(p, "{{ height }}") {
				enc[i] = strings.Replace(p, "{{ height }}", fmt.Sprintf("%d", height), -1)
			}
		}
	}
	for i, p := range fallback {
		if strings.Contains(p, "{{ input }}") {
			fallback[i] = strings.Replace(p, "{{ input }}", file+".y4m", -1)
		}
		if strings.Contains(p, "{{ output }}") {
			fallback[i] = strings.Replace(p, "{{ output }}", file+".ivf", -1)
		}
		if strings.Contains(p, "{{ width }}") || strings.Contains(p, "{{ height }}") {
			width, height, err := Dimension(log, file)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			if strings.Contains(p, "{{ width }}") {
				enc[i] = strings.Replace(p, "{{ width }}", fmt.Sprintf("%d", width), -1)
			}
			if strings.Contains(p, "{{ height }}") {
				enc[i] = strings.Replace(p, "{{ height }}", fmt.Sprintf("%d", height), -1)
			}
		}
	}
	for i, p := range repack {
		if strings.Contains(p, "{{ input }}") {
			repack[i] = strings.Replace(p, "{{ input }}", file+".ivf", -1)
		}
		if strings.Contains(p, "{{ output }}") {
			repack[i] = strings.Replace(p, "{{ output }}", without_ext+".avif", -1)
		}
	}
	return ext, enc, fallback, repack, nil
}
