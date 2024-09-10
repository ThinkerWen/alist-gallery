package handles

import (
	"alist-gallery/config"
	"alist-gallery/server/common"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
)

// FormGallery 以form形式上传图片到图床，对应alist中 /api/fs/form
func FormGallery(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	contentType := c.Request().Header.Get("Content-Type")
	filePath := c.Request().Header.Get("File-Path")
	fileName := c.Request().Header.Get("File-Name")
	asTask := c.Request().Header.Get("As-Task")

	if !strings.Contains(contentType, "multipart/form-data") {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Content-Type must be multipart/form-data"})
	}

	// 处理请求中的文件
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "No file part in the request"})
	}
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)

	headers := make(map[string]string)
	if filePath != "" {
		_, fileName = path.Split(filePath)
	} else if fileName != "" {
		filePath = path.Join(config.CONFIG.StoragePath, fileName)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "File name required"})
	}
	headers["As-Task"] = asTask
	headers["File-Path"] = filePath
	headers["Content-Type"] = contentType
	headers["Authorization"] = authorization

	client := resty.New()
	r, err := client.R().
		SetHeaders(headers).
		SetFileReader("file", file.Filename, f).
		SetFormData(map[string]string{"test": "1"}).
		Put(config.CONFIG.AListHost + "/api/fs/form")
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}

	data := string(r.Body())
	if gjson.Get(data, "code").Int() != 200 {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": "alist service error"})
	}
	resp := new(common.FsUP)
	if asTask == "true" {
		resp.Data.Task = json.RawMessage(gjson.Get(data, "data.task").String())
	}
	resp.Code = 200
	resp.Message = "success"
	resp.Data.Name = fileName
	resp.Data.Path = filePath
	resp.Data.Url = config.CONFIG.GalleryLocation + path.Join("/fs/show-gallery", fileName)
	return c.JSON(http.StatusOK, *resp)
}

// ShowImage 展示图床中图片（非下载）
func ShowImage(c echo.Context) error {
	imageName := c.Param("name")
	client := resty.New()
	r, err := client.R().
		SetHeader("Authorization", config.CONFIG.AListToken).
		SetBody(map[string]interface{}{"path": path.Join(config.CONFIG.StoragePath, imageName), "password": config.CONFIG.Password, "page": 1, "per_page": 1, "refresh": false}).
		Post(config.CONFIG.AListHost + "/api/fs/get")
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": err.Error()})
	}
	data := string(r.Body())
	if gjson.Get(data, "code").Int() != 200 {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": "alist service error"})
	}

	image, err := client.R().Get(gjson.Get(data, "data.raw_url").String())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{"message": "alist-gallery service error"})
	}
	return c.Blob(http.StatusOK, "image/jpeg", image.Body())
}
