package main

import (
	libs "batchavif/libs"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	var cf libs.Config
	cf.ParseConfig()

	var failedToConvert []string
	var original_files_size, converted_files_size float64 = 0, 0

	var log *os.File
	if cf.Config.ExportLog {
		log, _ = os.Create(fmt.Sprintf("%s.log", time.Now().Format("2006-01-02-15-04-05")))
		defer log.Close()
	}
	all_ext := append(cf.Image.Formats, cf.Animation.Formats...)
	files := libs.ListFiles(".", all_ext, cf.Config.Recursive)
	
	if !cf.Config.KeepOriginalExtension {
		if dupl_list := libs.ArrHasDupl(files); dupl_list != nil {
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
		for _, i := range files {
			name := i
			if !cf.Config.KeepOriginalExtension {
				name = strings.TrimSuffix(i, filepath.Ext(i))
			}
			if _, err := os.Stat(name); err == nil {
				os.Remove(name)
			}
		}
	}

	for _, file := range files {
		name := file
		if !cf.Config.KeepOriginalExtension {
			name = strings.TrimSuffix(file, filepath.Ext(file))
		}

		fmt.Println("==>", file)

		if _, err := os.Stat(name + ".avif"); os.IsExist(err) {
			fmt.Printf("Already converted\n\n")
			continue
		}

		var extractor []string
		var encoder []string
		var encoder_fallback []string
		var repackager []string

		if libs.InArr(cf.Image.Formats, filepath.Ext(file)) {
			extractor = cf.Image.Extractor
			encoder = cf.Image.Encoder
			encoder_fallback = cf.Image.EncoderFallback
			repackager = cf.Image.Repackager
		}

		if libs.InArr(cf.Animation.Formats, filepath.Ext(file)) {
			extractor = cf.Animation.Extractor
			encoder = cf.Animation.Encoder
			encoder_fallback = cf.Animation.EncoderFallback
			repackager = cf.Animation.Repackager
		}

		fmt.Println()
		err := libs.Convert(log, file, extractor, encoder, repackager, false)
		if err == nil {
			os.Remove(name + ".y4m")
			os.Remove(name + ".ivf")
			libs.ReportFileSize(libs.FileSize(file), libs.FileSize(name+".avif"))
			fmt.Println("")
			original_files_size += libs.FileSize(file)
			converted_files_size += libs.FileSize(name + ".avif")
			continue
		}

		fmt.Printf("ERROR: %s\n", err)

		if len(encoder_fallback) > 0 {
			os.Remove(name + ".ivf")
			err = libs.Convert(log, file, extractor, encoder_fallback, repackager, true)
			if err == nil {
				os.Remove(name + ".y4m")
				os.Remove(name + ".ivf")
				libs.ReportFileSize(libs.FileSize(file), libs.FileSize(name+".avif"))
				fmt.Println("")
				original_files_size += libs.FileSize(file)
				converted_files_size += libs.FileSize(name + ".avif")
				continue
			}
		}

		os.Remove(name + ".y4m")
		os.Remove(name + ".ivf")

		fmt.Printf("ERROR: %s\n\n", err)
		failedToConvert = append(failedToConvert, file)
	}

	if cf.Config.DeleteAfterConversion {
		for _, file := range files {
			if _, err := os.Stat(file); os.IsExist(err) {
				os.Remove(file)
			}
		}
	}

	if len(failedToConvert) > 0 {
		fmt.Printf("==> %d file(s) failed to convert\n", len(failedToConvert))
		for _, i := range failedToConvert {
			fmt.Println(i)
		}
		fmt.Println("")
	}

	fmt.Printf("==> %d file(s) converted\n", len(files)-len(failedToConvert))
	libs.ReportFileSize(original_files_size, converted_files_size)

	fmt.Printf("\nPress Enter to exit...")
	var exit string
	fmt.Scanln(&exit)
}
