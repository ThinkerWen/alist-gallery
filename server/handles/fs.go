package handles

import (
	"alist-gallery/config"
	"alist-gallery/internal/db"
	"alist-gallery/internal/model"
	"alist-gallery/server/common"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"net/http"
	"path"
)

// FormGallery 以form形式上传图片到图床，对应alist中 /api/fs/form
func FormGallery(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	fileName := c.Request().Header.Get("File-Name")
	asTask := c.Request().Header.Get("As-Task")

	user, err := common.GetUserName(authorization)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "User baned"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "No file part in the request"})
	}
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	if fileName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "File name required"})
	}
	r, err := common.FsFrom(fmt.Sprintf("%s/%s/%s", config.CONFIG.StoragePath, user, fileName), asTask, f)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}
	item := model.GalleryIndex{Path: config.CONFIG.StoragePath, User: user, ImageName: fileName}
	if err = db.SetIndexItem(item); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	resp := new(common.FsUP)
	if asTask == "true" {
		resp.Data.Task = json.RawMessage(gjson.Get(string(r), "data.task").String())
	}
	resp.Code = 200
	resp.Message = "success"
	resp.Data.Name = fileName
	resp.Data.Url = config.CONFIG.GalleryLocation + path.Join("/fs/show-gallery", fileName)
	return c.JSON(http.StatusOK, *resp)
}

// PutGallery 以stream形式上传图片到图床，对应alist中 /api/fs/put
func PutGallery(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	fileName := c.Request().Header.Get("File-Name")
	asTask := c.Request().Header.Get("As-Task")
	user, err := common.GetUserName(authorization)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "User baned"})
	}

	if fileName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "File name required"})
	}
	r, err := common.FsStream(path.Join(config.CONFIG.StoragePath, user, fileName), asTask, c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}
	item := model.GalleryIndex{Path: config.CONFIG.StoragePath, User: user, ImageName: fileName}
	if err = db.SetIndexItem(item); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	resp := new(common.FsUP)
	if asTask == "true" {
		resp.Data.Task = json.RawMessage(gjson.Get(string(r), "data.task").String())
	}
	resp.Code = 200
	resp.Message = "success"
	resp.Data.Name = fileName
	resp.Data.Url = config.CONFIG.GalleryLocation + path.Join("/fs/show-gallery", fileName)
	return c.JSON(http.StatusOK, *resp)
}
