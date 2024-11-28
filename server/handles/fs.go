package handles

import (
	"alist-gallery/config"
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "User baned"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "No file part in the request"})
	}
	if fileName == common.Blank {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "File name required"})
	}
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := common.FsFrom(fmt.Sprintf(common.StorageFormatter, user, fileName), asTask, f)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	resp := new(common.FsUP)
	if asTask == "true" {
		resp.Data.Task = json.RawMessage(gjson.Get(string(r), "data.task").String())
	}
	go common.FsRefresh()
	resp.Code = 200
	resp.Message = "success"
	resp.Data.Name = fileName
	resp.Data.Url = common.GalleryFormatter + fileName
	return c.JSON(http.StatusOK, *resp)
}

// PutGallery 以stream形式上传图片到图床，对应alist中 /api/fs/put
func PutGallery(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	fileName := c.Request().Header.Get("File-Name")
	asTask := c.Request().Header.Get("As-Task")
	user, err := common.GetUserName(authorization)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "User baned"})
	}

	if fileName == common.Blank {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "File name required"})
	}

	r, err := common.FsStream(path.Join(config.CONFIG.StoragePath, user, fileName), asTask, c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	resp := new(common.FsUP)
	if asTask == "true" {
		resp.Data.Task = json.RawMessage(gjson.Get(string(r), "data.task").String())
	}
	go common.FsRefresh()
	resp.Code = 200
	resp.Message = "success"
	resp.Data.Name = fileName
	resp.Data.Url = common.GalleryFormatter + fileName
	return c.JSON(http.StatusOK, *resp)
}
