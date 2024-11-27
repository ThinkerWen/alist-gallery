package common

import (
	"alist-gallery/config"
	"net/url"
	"path"
)

var (
	Blank = ""

	ErrNoSuchItem = "no such item"

	StorageFormatter = path.Join(config.CONFIG.StoragePath, "/%s/%s")

	GalleryFormatter, _ = url.JoinPath(config.CONFIG.GalleryLocation, "/fs/show-gallery", "%s")

	ApiMe, _ = url.JoinPath(config.CONFIG.AListHost, "/api/me")

	ApiFsGet, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/get")

	ApiFsPut, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/put")

	ApiFsForm, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/form")

	ApiFsSearch, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/search")
)
