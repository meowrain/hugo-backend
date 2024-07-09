package config

import (
	"embed"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Auth struct {
		UserName string `yaml:"username"`
		PassWord string `yaml:"password"`
	} `yaml:"auth"`
}
type GlobalConfig struct {
	UserName string
	PassWord string
}

//go:embed config.yaml
var EmbeddedConfig embed.FS

var cfg Config
var GlobalConfigInstance GlobalConfig

func init() {
	GlobalConfigInstance = GlobalConfig{}
	err := LoadFromFile()
	if err != nil {
		log.Println(err)
	}
}

func LoadFromFile() error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		data, err = EmbeddedConfig.ReadFile("config.yaml")
		if err != nil {
			return err
		}
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	GlobalConfigInstance.UserName = cfg.Auth.UserName
	GlobalConfigInstance.PassWord = cfg.Auth.PassWord
	return nil
}
