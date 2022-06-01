package libs

import (
	"bytes"
	"fmt"
	pipe "github.com/b4b4r07/go-pipe"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// run extract-encode-repack file by file, export log, print detailed progress
func StandardConvert(
	log *os.File,
	file string,
	ext []string,
	enc []string,
	repack []string,
	rerun bool,
) error {
	if !rerun {
		logDivider(log, "EXTRACT TO Y4M", ext[0], ext[1:])
		fmt.Println("Extracting to y4m format...")
		if ExecCommand(log, ext[0], ext[1:]...) != nil {
			return fmt.Errorf("failed to extract")
		}
		logDivider(log, "CONVERT TO IVF", enc[0], enc[1:])
		fmt.Printf("Converting to avif using %s...\n", enc[0])
		if ExecCommand(log, enc[0], enc[1:]...) != nil {
			return fmt.Errorf("failed to convert")
		}
	}
	if rerun {
		logDivider(log, "RETRY CONVERT TO IVF", enc[0], enc[1:])
		fmt.Printf("Retryng with %s...\n", enc)
		if ExecCommand(log, enc[0], enc[1:]...) != nil {
			return fmt.Errorf("failed to convert")
		}
	}
	logDivider(log, "REPACK TO AVIF", repack[0], repack[1:])
	fmt.Println("Repacking to avif...")
	if ExecCommand(log, repack[0], repack[1:]...) != nil {
		return fmt.Errorf("failed to repack")
	}
	return nil
}

// Like ConvertFileModeSingleThread, but without logging or printing progress
func spawnFileModeJob(log *os.File, file string, ext []string, enc []string, repack []string, rerun bool) error {
	if !rerun {
		if ExecCommand(nil, ext[0], ext[1:]...) != nil {
			return fmt.Errorf("failed to extract")
		}
		if ExecCommand(nil, enc[0], enc[1:]...) != nil {
			return fmt.Errorf("failed to convert")
		}
	}
	if rerun {
		if ExecCommand(nil, enc[0], enc[1:]...) != nil {
			return fmt.Errorf("failed to convert")
		}
	}
	if ExecCommand(nil, repack[0], repack[1:]...) != nil {
		return fmt.Errorf("failed to repack")
	}
	return nil
}

// like spawnFileModeJob, but piping commands directly without extracting middle files
func spawnPipeModeJob(ext []string, enc []string, repack []string) error {
	var b bytes.Buffer
	err := pipe.Command(&b,
		exec.Command(ext[0], ext[1:]...),
		exec.Command(enc[0], enc[1:]...),
		exec.Command(repack[0], repack[1:]...),
	)
	if err != nil {
		return err
	}
	return nil
}

// spawing multiple spawnFileModeJob/spawnPipeModeJob in parallel
func Convert(
	log *os.File,
	files chan string,
	ext_ []string,
	enc_ []string,
	fallback_ []string,
	repack_ []string,
	mode string,
	fail_convert *[]string,
	skip_convert *[]string,
	orig_sizes *float64,
	converted_sizes *float64,
	wg *sync.WaitGroup,
) {
	var cf Config
	cf.ParseConfig()
	for file := range files {
		start := time.Now()
		// region: FILE ALREADY EXIST?
		name := file
		if !cf.Config.KeepOriginalExtension {
			name = strings.TrimSuffix(file, filepath.Ext(file))
		}
		// if name + ".avif" is exist, add file to skip_convert and continue
		if _, err := os.Stat(name + ".avif"); err == nil {
			fmt.Printf("==> %s.avif already existed\n", name)
			*skip_convert = append(*skip_convert, file)
			continue
		}
		// endregion

		// CREATE COMMANDS
		ext, enc, fallback, repack, err := ProcessPreset(nil, file, ext_, enc_, fallback_, repack_)
		if err != nil {
			*fail_convert = append(*fail_convert, file)
			fmt.Println("==> ERROR [config validation]:", file, "\n", err)
			continue
		}

		// START CONVERT
		var errMain error
		if mode == "file" {
			errMain = spawnFileModeJob(log, file, ext, enc, repack, false)
		} else if mode == "pipe" {
			errMain = spawnPipeModeJob(ext, enc, repack)
		}
		os.Remove(file + ".ivf") // remove ivf created by failed encoder, keep the extracted y4m for fallback

		// if failed, no fallback
		if (errMain != nil) && (len(fallback) <= 0) {
			os.Remove(file + ".y4m") // remove y4m too if no fallback
			*fail_convert = append(*fail_convert, file)
			fmt.Println("==> ERROR [main encoder]:", file, "\n", errMain)
			continue
		} else if len(fallback) > 0 {
			// fallback
			var errFallback error
			if mode == "file" {
				errFallback = spawnFileModeJob(log, file, ext, fallback, repack, true)
			}
			if mode == "pipe" {
				errFallback = spawnPipeModeJob(ext, fallback, repack)
			}
			os.Remove(file + ".y4m")
			os.Remove(file + ".ivf")
			// if fallback also failed
			if errFallback != nil {
				*fail_convert = append(*fail_convert, file)
				fmt.Println("==> ERROR [fallback encoder]:", file, "\n", errFallback)
				continue
			}
		}
		// if success
		os.Remove(file + ".y4m") // remove y4m if success since we've already removed the ivf right after errMain
		*orig_sizes += FileSize(file)
		*converted_sizes += FileSize(name + ".avif")
		// rount time since start to 2 decimal places

		fmt.Printf("==> SUCCESS: %s | %s | %s\n", file, ReportFileSize(FileSize(file), FileSize(name+".avif")), Timer(&start))
	}
	wg.Done()
}
