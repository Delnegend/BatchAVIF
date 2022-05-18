package libs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Replace {{ input }}, {{ output }},... with actual values
func processInput(
	log *os.File,
	input string,
	ext []string,
	enc []string,
	repack []string,
) (string, []string, string, []string, string, []string, error) {
	var cf Config
	cf.ParseConfig()

	without_ext := input
	if !cf.Config.KeepOriginalExtension {
		without_ext = strings.TrimSuffix(input, filepath.Ext(input))
	}

	extractor := ext[0]
	extract_preset := make([]string, len(ext)-1)
	copy(extract_preset, ext[1:])
	encoder := enc[0]
	encoder_preset := make([]string, len(enc)-1)
	copy(encoder_preset, enc[1:])
	repackager := repack[0]
	repack_preset := make([]string, len(repack)-1)
	copy(repack_preset, repack[1:])

	for i, p := range extract_preset {
		if strings.Contains(p, "{{ input }}") {
			extract_preset[i] = input
		}
		if strings.Contains(p, "{{ output }}") {
			extract_preset[i] = without_ext + ".y4m"
		}
	}
	if encoder == "aomenc" {
		width, height, err := Dimension(log, input)
		if err != nil {
			return "", nil, "", nil, "", nil, err
		}
		for i, p := range encoder_preset {
			if strings.Contains(p, "{{ width }}") {
				encoder_preset[i] = strings.Replace(p, "{{ width }}", fmt.Sprintf("%d", width), -1)
			}
			if strings.Contains(p, "{{ height }}") {
				encoder_preset[i] = strings.Replace(p, "{{ height }}", fmt.Sprintf("%d", height), -1)
			}
		}
	}
	for i, p := range encoder_preset {
		if strings.Contains(p, "{{ input }}") {
			encoder_preset[i] = strings.Replace(p, "{{ input }}", without_ext+".y4m", -1)
		}
		if strings.Contains(p, "{{ output }}") {
			encoder_preset[i] = strings.Replace(p, "{{ output }}", without_ext+".ivf", -1)
		}
		if strings.Contains(p, "{{ threads }}") {
			encoder_preset[i] = strings.Replace(p, "{{ threads }}", fmt.Sprintf("%d", MaxCPU()), -1)
		}
	}
	for i, p := range repack_preset {
		if strings.Contains(p, "{{ input }}") {
			repack_preset[i] = strings.Replace(p, "{{ input }}", without_ext+".ivf", -1)
		}
		if strings.Contains(p, "{{ output }}") {
			repack_preset[i] = strings.Replace(p, "{{ output }}", without_ext+".avif", -1)
		}
	}
	return extractor, extract_preset, encoder, encoder_preset, repackager, repack_preset, nil
}

// Just a fancy way to divide sections in the log
func logDivider(log *os.File, message string, enc string, preset []string) {
	fmt.Fprintf(log, "\n\n\n==================== %s ====================\n", message)
	fmt.Fprintf(log, "Using %s with preset: %s\n", enc, preset)
	fmt.Fprintf(log, "\n\n\n")
}

func Convert(
	log *os.File,
	file string,
	ext []string,
	enc []string,
	repack []string,
	rerun bool,
) error {

	extractor, extract_preset, encoder, encoder_preset, repackager, repack_preset, err := processInput(log, file, ext, enc, repack)
	if err != nil {
		return err
	}

	if !rerun {
		logDivider(log, "EXTRACT TO Y4M", extractor, extract_preset)
		fmt.Println("Extracting to y4m format...")
		if ExecCommand(log, extractor, extract_preset...) != nil {
			return errors.New("failed to extract")
		}

		logDivider(log, "CONVERT TO IVF", encoder, encoder_preset)
		fmt.Printf("Converting to avif using %s...\n", encoder)
		if ExecCommand(log, encoder, encoder_preset...) != nil {
			return errors.New("failed to convert")
		}
	}
	if rerun {
		logDivider(log, "RETRY CONVERT TO IVF", encoder, encoder_preset)
		fmt.Printf("Retryng with %s...\n", encoder)
		if ExecCommand(log, encoder, encoder_preset...) != nil {
			return errors.New("failed to convert")
		}
	}

	logDivider(log, "REPACK TO AVIF", repackager, repack_preset)
	fmt.Printf("Repacking to avif using %s...\n", repackager)
	if ExecCommand(log, repackager, repack_preset...) != nil {
		return errors.New("failed to repack")
	}

	return nil
}
