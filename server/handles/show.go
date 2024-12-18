package handles

import (
	"alist-gallery/config"
	"alist-gallery/internal/db"
	"alist-gallery/internal/model"
	"alist-gallery/internal/net"
	"alist-gallery/server/common"
	"bytes"
	"errors"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"path"
	"strconv"
	"time"
)

// ShowImage 展示图床中图片（非下载）
func ShowImage(c echo.Context) error {
	name := c.Param("name")
	qualityStr := c.QueryParam("size")
	quality := config.CONFIG.Compression
	if qualityStr != "" && quality != 0 {
		quality, _ = strconv.Atoi(qualityStr)
	}

	cacheKey := fmt.Sprintf(common.RedisFormatter, name)
	cachedImage := db.RedisGet(cacheKey)
	if cachedImage != "" {
		c.Response().Header().Set("Cache-Control", "public, max-age=2592000")
		img, contentType, err := compressImageToWebP([]byte(cachedImage), quality)
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
		}
		return c.Blob(http.StatusOK, contentType, img)
	}

	item, err := db.GetGalleryItem(name)
	if err != nil || item.ImageURL == "" {
		item, err = searchImage(name)
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
		}
		if err = db.SetGalleryItem(item); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
		}
	}

	img, err := loadImage(item.ImageURL)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	if len(img) > 1024 {
		db.RedisSet(cacheKey, string(img), time.Duration(config.CONFIG.Redis.Timeout)*time.Minute)
	}
	img, contentType, err := compressImageToWebP([]byte(cachedImage), quality)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}
	c.Response().Header().Set("Cache-Control", "public, max-age=2592000")
	return c.Blob(http.StatusOK, contentType, img)
}

func searchImage(name string) (model.GalleryIndex, error) {
	res, err := common.FsSearch(name)
	if err != nil {
		return model.GalleryIndex{}, err
	}
	data := gjson.Get(string(res), "data.content|0")
	imageLink, err := common.FsGet(path.Join(data.Get("parent").String(), data.Get("name").String()))
	if imageLink == common.Blank || err != nil {
		return model.GalleryIndex{}, errors.New("alist service image not found")
	}

	return model.GalleryIndex{
		Path:      data.Get("parent").String(),
		User:      path.Base(data.Get("parent").String()),
		ImageName: name,
		ImageURL:  imageLink,
	}, nil
}

func loadImage(imageLink string) ([]byte, error) {
	resp, err := net.GlobalClient.R().Get(imageLink)
	if err != nil {
		return nil, errors.New("alist-gallery service error")
	}
	return resp.Body(), nil
}

func compressImageToWebP(data []byte, quality int) ([]byte, string, error) {
	if quality == 0 {
		return data, "image/png", nil
	}
	if quality > 100 || quality < 0 {
		quality = 100
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}
	if format == "gif" {
		return data, "image/gif", nil
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)}); err != nil {
		return nil, "", fmt.Errorf("failed to encode WebP image: %w", err)
	}

	return buf.Bytes(), "image/webp", nil
}
