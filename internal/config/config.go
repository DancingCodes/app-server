package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

var GlobalConfig Config

// LoadConfig 从文件加载配置
func LoadConfig(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	err = yaml.Unmarshal(file, &GlobalConfig)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}
