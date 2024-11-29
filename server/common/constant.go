package common

import (
	"alist-gallery/config"
	"net/url"
	"path"
)

var (
	Blank = ""

	ErrNoSuchItem = "no such item"

	RedisFormatter = "alist-gallery:%s"

	StorageFormatter = path.Join(config.CONFIG.StoragePath, "/%s/%s")

	GalleryFormatter, _ = url.JoinPath(config.CONFIG.GalleryLocation, "/fs/show-gallery/")

	ApiMe, _ = url.JoinPath(config.CONFIG.AListHost, "/api/me")

	ApiFsGet, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/get")

	ApiFsPut, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/put")

	ApiFsForm, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/form")

	ApiFsSearch, _ = url.JoinPath(config.CONFIG.AListHost, "/api/fs/search")

	ApiIndexUpdate, _ = url.JoinPath(config.CONFIG.AListHost, "/api/admin/index/update")
)
