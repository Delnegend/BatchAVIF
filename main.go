package main

import (
	libs "batchavif/libs"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func reportFileSize(old int64, new int64) {
	fmt.Printf(
		"Original: %s | Converted: %s ~ %.2f%%\n",
		libs.HumanReadableSize(old),
		libs.HumanReadableSize(new),
		float64(new)/float64(old)*100,
	)
}

func main() {

	var cf libs.Config
	cf.ParseConfig()

	var failedToConvert []string

	var original_files_size, converted_files_size int64 = 0, 0

	var log *os.File
	if cf.Config.ExportLog {
		log, _ = os.Create(fmt.Sprintf("%s.log", time.Now().Format("2006-01-02-15-04-05")))
		defer log.Close()
	}

	images := libs.ListFiles(".", cf.Image.Formats, cf.Config.Recursive)
	animations := libs.ListFiles(".", cf.Animation.Formats, cf.Config.Recursive)

	all := append(images, animations...)

	if !cf.Config.KeepOriginalExtension {
		dupl, dupl_list := libs.HasDuplName(all)
		if dupl {
			fmt.Println("ERROR: Duplicate file names found. Please rename them or set keep_original_extension to true.")
			for _, i := range dupl_list {
				fmt.Println(i)
			}
			os.Exit(1)
		}
	}

	if cf.Config.DeleteAfterConversion {
		fmt.Println("WARNING: Keep original files is off. Confirm delete original files after conversion? [y/N].")
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
		} else {
			for _, i := range all {
				name := i
				if !cf.Config.KeepOriginalExtension {
					name = strings.TrimSuffix(i, filepath.Ext(i))
				}
				if _, err := os.Stat(name); err == nil {
					os.Remove(name)
				}
			}
		}
	}

	for _, image := range images {
		name := image
		if !cf.Config.KeepOriginalExtension {
			name = strings.TrimSuffix(image, filepath.Ext(image))
		}

		fmt.Println("==>", name)

		if _, err := os.Stat(name + ".avif"); os.IsExist(err) {
			fmt.Printf("Already converted\n\n")
			continue
		}

		err := libs.ConvertImg(log, image, cf.Image.Extractor, cf.Image.Encoder, cf.Image.Repackager)
		if err == nil {
			reportFileSize(libs.FileSize(image), libs.FileSize(name+".avif"))
			fmt.Println("")
			original_files_size += libs.FileSize(image)
			converted_files_size += libs.FileSize(name + ".avif")
			continue
		}

		fmt.Printf("ERROR: %s\n\n", err)
		failedToConvert = append(failedToConvert, image)
	}

	for _, ani := range animations {
		name := ani
		if !cf.Config.KeepOriginalExtension {
			name = strings.TrimSuffix(ani, filepath.Ext(ani))
		}

		fmt.Println("==>", ani)

		if _, err := os.Stat(name + ".avif"); os.IsExist(err) {
			fmt.Printf("Already converted\n\n")
			continue
		}

		err := libs.ConvertAni(log, ani, cf.Animation.Extractor, cf.Animation.EncoderMain, cf.Animation.Repackager, false)
		if err == nil {
			reportFileSize(libs.FileSize(ani), libs.FileSize(name+".avif"))
			fmt.Println("")
			original_files_size += libs.FileSize(ani)
			converted_files_size += libs.FileSize(name + ".avif")
			continue
		}
		fmt.Printf("ERROR: %s\n", err)

		err2 := libs.ConvertAni(log, ani, cf.Animation.Extractor, cf.Animation.EncoderFallback, cf.Animation.Repackager, true)
		if err2 == nil {
			reportFileSize(libs.FileSize(ani), libs.FileSize(name+".avif"))
			fmt.Println("")
			original_files_size += libs.FileSize(ani)
			converted_files_size += libs.FileSize(name + ".avif")
			continue
		}

		fmt.Printf("ERROR: %s\n\n", err2)
		failedToConvert = append(failedToConvert, ani)
	}

	if cf.Config.DeleteAfterConversion {
		for _, file := range all {
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

	fmt.Printf("==> %d file(s) converted\n", len(images)+len(animations)-len(failedToConvert))
	reportFileSize(original_files_size, converted_files_size)
}
