package common

import (
	"alist-gallery/config"
	"path"
)

var (
	Blank = ""

	ErrNoSuchItem = "no such item"

	StorageFormatter = path.Join(config.CONFIG.StoragePath, "/%s/%s")

	GalleryFormatter = path.Join(config.CONFIG.GalleryLocation, "/fs/show-gallery", "%s")

	ApiMe = path.Join(config.CONFIG.AListHost, "/api/me")

	ApiFsGet = path.Join(config.CONFIG.AListHost, "/api/fs/get")

	ApiFsPut = path.Join(config.CONFIG.AListHost, "/api/fs/put")

	ApiFsForm = path.Join(config.CONFIG.AListHost, "/api/fs/form")

	ApiFsSearch = path.Join(config.CONFIG.AListHost, "/api/fs/search")
)
