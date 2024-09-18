package model

import "time"

type GalleryIndex struct {
	Id        string    `json:"id"`
	Path      string    `json:"path"`
	User      string    `json:"user"`
	ImageName string    `json:"image_name"`
	CreatedAt time.Time `json:"created_at"`
}
