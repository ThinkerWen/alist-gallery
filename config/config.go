package config

import (
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type Config struct {
	Port            int    `yaml:"port"`
	AListHost       string `yaml:"alist-host"`
	GalleryLocation string `yaml:"gallery-location"`
	StoragePath     string `yaml:"storage-path"`
	AListToken      string `yaml:"alist-token"`
	Password        string `yaml:"password"`
}

var CONFIG Config
var mu sync.Mutex

func init() {
	workDir, _ := os.Getwd()
	yamlFile, err := os.ReadFile(workDir + "/config.yaml")
	if err != nil {
		initDefaultConfig()
		if err = SaveConfig(); err != nil {
			log.Fatalf("Error saving config: %v", err)
		}
	} else if err = yaml.Unmarshal(yamlFile, &CONFIG); err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	log.Info("Load config successfully")
}

func initDefaultConfig() {
	config := new(Config)
	config.Port = 5243
	config.AListHost = "https://assets.example.com"
	config.StoragePath = "/Storage/Gallery"
	config.GalleryLocation = "https://assets.example.com:5243"
	config.AListToken = "alist-4254afdc-1acg-1999-08aa-b6134kx4kv63FdkHJFPeaFDdEGYmSe29KETy4fdsareKM8fdsagfdsgfdgfdagdfgr"
	config.Password = ""

	CONFIG = *config
	log.Info("Set default config successfully")
}

// SaveConfig 将配置保存回文件
func SaveConfig() error {
	mu.Lock()
	defer mu.Unlock()
	workDir, _ := os.Getwd()
	yamlData, err := yaml.Marshal(CONFIG)
	if err != nil {
		return err
	}

	if err = os.WriteFile(workDir+"/config.yaml", yamlData, 0644); err != nil {
		return err
	}

	log.Info("Configuration saved successfully")
	return nil
}
