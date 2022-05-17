package libs

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func usingLibs(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	img, _, err2 := image.DecodeConfig(file)
	if err2 != nil {
		return 0, 0, err2
	}
	return img.Width, img.Height, nil
}

func usingFfmpeg(log *os.File, path string) (int, int, error) {
	ffmpeg_extract_frame := []string{"-i", path, "-vframes", "1", "-f", "image2", path+".frame.png"}
	logDivider(log, "GETTING DIMENSION", "ffmpeg", ffmpeg_extract_frame)
	ExecCommand(log, "ffmpeg", ffmpeg_extract_frame...)
	width, height, err := usingLibs(path + ".frame.png")
	errRemove := os.Remove(fmt.Sprintf("%s.frame.png", path))
	if errRemove != nil {
		fmt.Printf("WARNING: %s.frame.png used for getting dimension could not be removed", path)
	}
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func Dimension(log *os.File, path string) (int, int, error) {
	if InArr([]string{".jpg", ".jpeg", ".png", ".gif"}, filepath.Ext(path)) {
		if width, height, err := usingLibs(path); err == nil {
			return width, height, nil
		}
		if width, height, err := usingFfmpeg(log, path); err == nil {
			return width, height, nil
		}
	}
	if width, height, err := usingFfmpeg(log, path); err == nil {
		return width, height, nil
	}
	return 0, 0, errors.New(path + " is not a valid image or FFMPEG cannot decode it")
}
