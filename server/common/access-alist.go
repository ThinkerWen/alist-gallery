package common

import (
	"alist-gallery/config"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"io"
	"mime/multipart"
)

func GetUserName(authorization string) (string, error) {
	client := resty.New()
	r, err := client.R().
		SetHeader("Authorization", authorization).
		Get(config.CONFIG.AListHost + "/api/me")
	if err != nil {
		return "", err
	}

	body := string(r.Body())
	if gjson.Get(body, "code").Int() != 200 {
		return "", errors.New("alist /api/me access error")
	}
	if gjson.Get(body, "disabled").Bool() {
		return "", errors.New("user has been disabled in alist")
	}
	return gjson.Get(body, "data.username").String(), nil
}

func FsGet(path string) (string, error) {
	client := resty.New()
	r, err := client.R().
		SetHeader("Authorization", config.CONFIG.AListToken).
		SetBody(map[string]interface{}{"path": path, "password": config.CONFIG.Password, "page": 1, "per_page": 1, "refresh": false}).
		Post(config.CONFIG.AListHost + "/api/fs/get")
	if err != nil {
		return "", err
	}

	body := string(r.Body())
	if gjson.Get(body, "code").Int() != 200 {
		return "", errors.New("alist /api/fs/get access error")
	}
	return gjson.Get(body, "data.raw_url").String(), nil
}

func FsFrom(filePath, asTask string, f multipart.File) ([]byte, error) {
	headers := make(map[string]string)
	headers["As-Task"] = asTask
	headers["File-Path"] = filePath
	headers["Content-Type"] = "multipart/form-data;"
	headers["Authorization"] = config.CONFIG.AListToken
	client := resty.New()
	r, err := client.R().
		SetHeaders(headers).
		SetFileReader("file", "name.png", f).
		Put(config.CONFIG.AListHost + "/api/fs/form")
	if err != nil {
		return nil, err
	}
	if gjson.Get(string(r.Body()), "code").Int() != 200 {
		return nil, errors.New("alist /api/fs/form access error")
	}
	return r.Body(), nil
}

func FsStream(filePath, asTask string, f io.ReadCloser) ([]byte, error) {
	headers := make(map[string]string)
	headers["As-Task"] = asTask
	headers["File-Path"] = filePath
	headers["Authorization"] = config.CONFIG.AListToken
	client := resty.New()
	r, err := client.R().
		SetHeaders(headers).
		SetBody(f).
		Put(config.CONFIG.AListHost + "/api/fs/put")
	if err != nil {
		return nil, err
	}
	if gjson.Get(string(r.Body()), "code").Int() != 200 {
		return nil, errors.New("alist /api/fs/put access error")
	}
	return r.Body(), nil
}
