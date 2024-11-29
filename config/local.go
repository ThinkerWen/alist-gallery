package config

import (
	"database/sql"
	"github.com/charmbracelet/log"
)

func initSqLite() {
	var err error
	if DB, err = sql.Open("sqlite3", "gallery.db"); err != nil {
		log.Fatalf("Error creating database: %v", err)
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
		log.Fatalf("Error creating database: %v", err)
	}
	log.Info("SqLite loaded")
}
