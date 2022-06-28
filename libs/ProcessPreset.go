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
	mode := cf.Config.Mode

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
		if strings.Contains(p, "{{ output }}") && (mode == "file") {
			ext[i] = strings.Replace(p, "{{ output }}", file+".y4m", -1)
		} else if mode == "pipe" {
			ext[i] = strings.Replace(p, "{{ output }}", "-", -1)
		}
	}
	for i, p := range enc {
		if strings.Contains(p, "{{ input }}") && (mode == "file") {
			enc[i] = strings.Replace(p, "{{ input }}", file+".y4m", -1)
		} else if mode == "pipe" {
			enc[i] = strings.Replace(p, "{{ input }}", "-", -1)
		}
		if strings.Contains(p, "{{ output }}") && (mode == "file") {
			enc[i] = strings.Replace(p, "{{ output }}", file+".ivf", -1)
		} else if mode == "pipe" {
			enc[i] = strings.Replace(p, "{{ output }}", "-", -1)
		}
		if strings.Contains(p, "{{ threads }}") {
			enc[i] = strings.Replace(p, "{{ threads }}", fmt.Sprintf("%d", MaxCPU()), -1)
		}
		// Moved {{ width }} and {{ height }} parsing to Convert.go for parsing dimension from y4m after extracted
	}
	for i, p := range fallback {
		if strings.Contains(p, "{{ input }}") && (mode == "file") {
			fallback[i] = strings.Replace(p, "{{ input }}", file+".y4m", -1)
		} else if mode == "pipe" {
			fallback[i] = strings.Replace(p, "{{ input }}", "-", -1)
		}
		if strings.Contains(p, "{{ output }}") && (mode == "file") {
			fallback[i] = strings.Replace(p, "{{ output }}", file+".ivf", -1)
		} else if mode == "pipe" {
			fallback[i] = strings.Replace(p, "{{ output }}", "-", -1)
		}
	}
	for i, p := range repack {
		if strings.Contains(p, "{{ input }}") && (mode == "file") {
			repack[i] = strings.Replace(p, "{{ input }}", file+".ivf", -1)
		} else if mode == "pipe" {
			repack[i] = strings.Replace(p, "{{ input }}", "-", -1)
		}
		if strings.Contains(p, "{{ output }}") {
			repack[i] = strings.Replace(p, "{{ output }}", without_ext+".avif", -1)
		}
	}
	return ext, enc, fallback, repack, nil
}
