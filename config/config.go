package config

import (
	"database/sql"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type Config struct {
	Port            int         `yaml:"port"`
	AListHost       string      `yaml:"alist-host"`
	GalleryLocation string      `yaml:"gallery-location"`
	StoragePath     string      `yaml:"storage-path"`
	AListToken      string      `yaml:"alist-token"`
	Password        string      `yaml:"password"`
	Compression     int         `yaml:"compression"`
	Redis           RedisConfig `yaml:"redis"`
}
type RedisConfig struct {
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
	Timeout  int    `yaml:"timeout"`
}

var DB *sql.DB
var RDB *redis.Client
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
	initRedis()
	initSqLite()
}

func initDefaultConfig() {
	CONFIG = Config{
		Port:            5243,
		AListHost:       "https://assets.example.com",
		GalleryLocation: "/Storage/Gallery",
		StoragePath:     "https://assets.example.com:5243",
		AListToken:      "alist-4254afdc-1acg-1999-08aa-...",
		Password:        "can be empty",
		Compression:     0,
		Redis: RedisConfig{
			Enable:   false,
			Host:     "127.0.0.1",
			Port:     6379,
			Database: 0,
			Password: "can be empty",
			Timeout:  60,
		},
	}
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
