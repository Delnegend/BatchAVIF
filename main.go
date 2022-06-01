package main

import (
	libs "batchavif/libs"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	// "math/rand"
)

func startConvert(
	log *os.File,
	files []string,
	ext_raw []string,
	enc_raw []string,
	fallback_raw []string,
	repack_raw []string,
	threads int,
	keepOriginalExtension bool,
	fail_convert *[]string,
	skip_convert *[]string,
	orig_sizes *float64,
	converted_sizes *float64,
) {
	var cf libs.Config
	cf.ParseConfig()
	mode, err := libs.ValidateConfig(cf.Config.Mode, ext_raw, enc_raw, fallback_raw, repack_raw)
	if threads <= 0 {
		threads = 1
	}
	if err != nil {
		panic(err)
	}

	// STANDARD CONVERSION, FILE BY FILE, COMMAND BY COMMAND, LOGGED
	if mode == "file" && threads == 1 {
		for _, file := range files {
			start := time.Now()
			fmt.Println("==>", file)
			// region: FILE ALREADY EXIST?
			name := file
			if !cf.Config.KeepOriginalExtension {
				name = strings.TrimSuffix(file, filepath.Ext(file))
			}
			if _, err := os.Stat(name + ".avif"); os.IsExist(err) {
				fmt.Printf("Already converted\n\n")
				*skip_convert = append(*skip_convert, file)
				continue
			}
			// endregion

			// region: CREATE COMMANDS
			ext, enc, fallback, rep, err := libs.ProcessPreset(log, file, ext_raw, enc_raw, fallback_raw, repack_raw)
			if err != nil {
				*fail_convert = append(*fail_convert, file)
				fmt.Println(err)
				continue
			}
			// endregion

			// region: START CONVERTING
			errMain := libs.StandardConvert(log, file, ext, enc, rep, false)
			os.Remove(file + ".ivf") // remove ivf created by failed encoder, keep the extracted y4m for fallback
			// if failed, no fallback
			if (errMain != nil) && len(fallback) <= 0 {
				os.Remove(file + ".y4m") // remove y4m too if no fallback
				fmt.Printf("ERROR: %s\n", err)
				*fail_convert = append(*fail_convert, file)
				continue
			} else if len(fallback) > 0 {
				// fallback
				errFallback := libs.StandardConvert(log, file, ext, fallback, enc, true)
				os.Remove(file + ".y4m")
				os.Remove(file + ".ivf")
				// if fallback also failed
				if errFallback != nil {
					fmt.Printf("ERROR: %s\n", err)
					*fail_convert = append(*fail_convert, file)
					continue
				}
			}

			// if success
			os.Remove(file + ".y4m") // remove y4m if success since we've already removed the ivf right after errMain
			*orig_sizes += libs.FileSize(file)
			*converted_sizes += libs.FileSize(name + ".avif")
			fmt.Printf("%s | %s", libs.ReportFileSize(libs.FileSize(file), libs.FileSize(name+".avif")), libs.Timer(&start))
			fmt.Println("")
			// endregion
		}
	}

	// PARALLEL CONVERSION, FILES BY FILES OR PIPPED COMMANDS, NO LOG.
	// It's like creating multiple standard conversion jobs, so checking file exist, creating commands and starting conversion is created in another func (libs.Convert) instead of here (like in standard mode above)
	if (threads > 1) || (mode == "pipe") {
		wg := new(sync.WaitGroup)
		files_queue := make(chan string)
		wg.Add(threads)

		for i := 1; i <= threads; i++ {
			go libs.Convert(nil, files_queue, ext_raw, enc_raw, fallback_raw, repack_raw, mode, fail_convert, skip_convert, orig_sizes, converted_sizes, wg)
		}
		for _, file := range files {
			files_queue <- file
		}
		close(files_queue)
		wg.Wait()
	}
}

func main() {

	// region: set variables, create log
	var cf libs.Config
	cf.ParseConfig()

	var fail_convert []string
	var skip_convert []string
	var original_sizes, converted_sizes float64 = 0, 0

	var log *os.File
	if cf.Config.ExportLog {
		log, _ = os.Create(fmt.Sprintf("%s.log", time.Now().Format("2006-01-02-15-04-05")))
		defer log.Close()
	}

	images := libs.ListFiles(".", cf.Image.Formats, cf.Config.Recursive)
	animations := libs.ListFiles(".", cf.Animation.Formats, cf.Config.Recursive)
	images_and_anis := append(images, animations...)
	// endregion

	// region: PRE-CONVERSION
	if !cf.Config.KeepOriginalExtension {
		if dupl_list := libs.ArrHasDupl(images_and_anis); dupl_list != nil {
			fmt.Println("ERROR: Duplicate file names found. Please rename them or set keep_original_extension to true.")
			for _, i := range dupl_list {
				fmt.Println(i)
			}
			os.Exit(1)
		}
	}
	if cf.Config.DeleteAfterConversion {
		fmt.Println("WARNING: Delete after conversion is on. Confirm delete original files after conversion? [y/N].")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			os.Exit(1)
		}
	}
	if cf.Config.Overwrite {
		fmt.Println("Overwrite mode is on. Continue? [y/N].")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			os.Exit(1)
		}
		for _, file := range images_and_anis {
			name := file
			if !cf.Config.KeepOriginalExtension {
				name = strings.TrimSuffix(file, filepath.Ext(file))
			}
			if _, err := os.Stat(name); err == nil {
				os.Remove(name)
			}
		}
	}
	// endregion

	start := time.Now()

	// images first
	startConvert(log, images, cf.Image.Extractor, cf.Image.Encoder, cf.Image.EncoderFallback, cf.Image.Repackager, cf.Config.Threads, cf.Config.KeepOriginalExtension, &fail_convert, &skip_convert, &original_sizes, &converted_sizes)
	// then animations
	startConvert(log, animations, cf.Animation.Extractor, cf.Animation.Encoder, cf.Animation.EncoderFallback, cf.Animation.Repackager, cf.Config.Threads, cf.Config.KeepOriginalExtension, &fail_convert, &skip_convert, &original_sizes, &converted_sizes)

	// region: POST CONVERSION
	// remove original files
	if cf.Config.DeleteAfterConversion {
		for _, file := range images_and_anis {
			if _, err := os.Stat(file); os.IsExist(err) {
				os.Remove(file)
			}
		}
	}
	// report converted
	fmt.Printf("\n==> %d file(s) converted\n", len(images_and_anis)-len(fail_convert)-len(skip_convert))
	fmt.Printf("%s | %s\n", libs.ReportFileSize(original_sizes, converted_sizes), libs.Timer(&start))
	fmt.Println()

	// report skipped
	if len(skip_convert) > 0 {
		fmt.Printf("==> %d file(s) skipped\n", len(skip_convert))
		for _, file := range skip_convert {
			fmt.Println(file)
		}
		fmt.Println()
	}

	// report failed conversion
	if len(fail_convert) > 0 {
		fmt.Printf("==> %d file(s) failed to convert\n", len(fail_convert))
		for _, i := range fail_convert {
			fmt.Println(i)
		}
		fmt.Println()
	}
	// endregion

	// EXIT
	fmt.Printf("Press Enter to exit...")
	var exit string
	fmt.Scanln(&exit)
}
