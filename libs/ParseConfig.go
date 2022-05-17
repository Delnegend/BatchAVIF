package libs

import (
	// "os"
	// "fmt"
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Example config config.yaml
	Image struct {
		Formats    []string `yaml:"formats"`
		Extractor  []string `yaml:"extractor"`
		Encoder    []string `yaml:"encoder"`
		Repackager []string `yaml:"repackager"`
	} `yaml:"image"`
	Animation struct {
		Formats         []string `yaml:"formats"`
		Extractor       []string `yaml:"extractor"`
		EncoderMain     []string `yaml:"encoder_main"`
		EncoderFallback []string `yaml:"encoder_fallback"`
		Repackager      []string `yaml:"repackager"`
	} `yaml:"animation"`
	Config struct {
		DeleteAfterConversion bool `yaml:"delete_after_conversion"`
		KeepOriginalExtension bool `yaml:"keep_original_extension"`
		Overwrite             bool `yaml:"overwrite"`
		Recursive             bool `yaml:"recursive"`
		ExportLog             bool `yaml:"export_log"`
	}
}

func configPath() (string, error) {
	if len(os.Args) > 1 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			return os.Args[1], nil
		}
	}
	if _, err := os.Stat("./config.yaml"); err == nil {
		return "./config.yaml", nil
	}
	return "", errors.New("config file not found")
}

func (config *Config) ParseConfig() *Config {
	path, err := configPath()
	if err != nil {
		panic(err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		panic(err)
	}
	return config
}
