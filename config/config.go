package config

import (
	"database/sql"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
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

var DB *sql.DB
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
	if err = createDatabase(); err != nil {
		log.Fatalf("Error creating database: %v", err)
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

func createDatabase() error {
	var err error
	if DB, err = sql.Open("sqlite3", "gallery.db"); err != nil {
		return err
	}

	createGalleryIndex := `
	CREATE TABLE IF NOT EXISTS gallery_index
	(
		id         INTEGER primary key autoincrement,
		path       VARCHAR(255)  default ''                not null,
		user       VARCHAR(255)  default ''                not null,
		image_name VARCHAR(255)  default ''                not null,
		image_url  varchar(2000) default ''                not null,
		created_at TIMESTAMP     default CURRENT_TIMESTAMP not null
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_image_name on gallery_index (image_name);`
	if _, err = DB.Exec(createGalleryIndex); err != nil {
		return err
	}

	return nil
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
