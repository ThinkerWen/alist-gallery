package handles

import (
	"alist-gallery/internal/db"
	"alist-gallery/internal/model"
	"alist-gallery/server/common"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

// ShowImage 展示图床中图片（非下载）
func ShowImage(c echo.Context) error {
	client := resty.New()
	name := c.Param("name")
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

	image, err := loadImage(item.ImageURL, client)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}
	return c.Blob(http.StatusOK, "image/png", image)
}

func searchImage(name string) (model.GalleryIndex, error) {
	res, err := common.FsSearch(name)
	if err != nil {
		return model.GalleryIndex{}, err
	}
	data := gjson.Get(string(res), "data.content|0")
	imageLink, err := common.FsGet(data.Get("parent").String())
	if err != nil {
		return model.GalleryIndex{}, err
	}

	return model.GalleryIndex{
		Path:      data.Get("parent").String(),
		User:      strings.Split(data.Get("parent").String(), "/")[1],
		ImageName: name,
		ImageURL:  imageLink,
	}, nil
}

func loadImage(imageLink string, client *resty.Client) ([]byte, error) {
	image, err := client.R().Get(imageLink)
	if err != nil {
		return nil, errors.New("alist-gallery service error")
	}
	return image.Body(), nil
}
