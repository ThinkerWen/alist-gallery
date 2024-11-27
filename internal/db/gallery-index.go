package db

import (
	"alist-gallery/config"
	"alist-gallery/internal/model"
	"alist-gallery/server/common"
	"errors"
)

func GetGalleryItem(name string) (model.GalleryIndex, error) {
	sqlStr := `SELECT * FROM gallery_index WHERE image_name = ? LIMIT 1`
	rows, err := config.DB.Query(sqlStr, name)
	if err != nil {
		return model.GalleryIndex{}, err
	}
	defer rows.Close()
	var item model.GalleryIndex
	if rows.Next() {
		if err := rows.Scan(&item.Id, &item.Path, &item.User, &item.ImageName, &item.ImageURL, &item.CreatedAt); err != nil {
			return model.GalleryIndex{}, err
		}
	} else {
		return model.GalleryIndex{}, errors.New(common.ErrNoSuchItem)
	}

	return item, nil
}

func SetGalleryItem(item model.GalleryIndex) error {
	sqlStr := `INSERT INTO gallery_index (path, user, image_name, image_url) VALUES (?, ?, ?)`
	if _, err := config.DB.Exec(sqlStr, item.Path, item.User, item.ImageName, item.ImageURL); err != nil {
		return err
	}
	return nil
}
