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
	if img, _, err := image.DecodeConfig(file); err == nil {
		return img.Width, img.Height, nil
	} else {
		return 0, 0, err
	}
}

func usingFfmpeg(log *os.File, path string) (int, int, error) {
	cmd := []string{"-i", path, "-vframes", "1", "-f", "image2", path + ".frame.png"}
	logDivider(log, "GETTING DIMENSION", "ffmpeg", cmd)
	if err := ExecCommand(log, "ffmpeg", cmd...); err != nil {
		return 0, 0, err
	}
	if width, height, err := usingLibs(path + ".frame.png"); err == nil {
		os.Remove(fmt.Sprintf("%s.frame.png", path))
		return width, height, err
	}
	return 0, 0, errors.New(path + " is not a valid image or FFMPEG cannot decode it to get dimension")
}

func Dimension(log *os.File, path string) (int, int, error) {
	if InArr([]string{".jpg", ".jpeg", ".png", ".gif"}, filepath.Ext(path)) {
		if width, height, err := usingLibs(path); err == nil {
			return width, height, err
		}
	}
	if width, height, err := usingFfmpeg(log, path); err == nil {
		return width, height, nil
	} else {
		return 0, 0, err
	}
}
