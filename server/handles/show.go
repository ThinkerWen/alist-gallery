package handles

import (
	"alist-gallery/config"
	"alist-gallery/internal/db"
	"alist-gallery/server/common"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"path"
)

// ShowImage 展示图床中图片（非下载）
func ShowImage(c echo.Context) error {
	name := c.Param("name")
	client := resty.New()
	item, err := db.GetIndexItem(name)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}
	imagePath := path.Join(config.CONFIG.StoragePath, item.User, name)
	imageLink, err := common.FsGet(imagePath)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	image, err := client.R().Get(imageLink)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": "alist-gallery service error"})
	}
	return c.Blob(http.StatusOK, "image/jpeg", image.Body())
}
