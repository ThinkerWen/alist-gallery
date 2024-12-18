package common

import (
	"alist-gallery/config"
	"alist-gallery/internal/net"
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"mime/multipart"
)

func GetUserName(authorization string) (string, error) {
	r, err := net.GlobalClient.R().
		SetHeader("Authorization", authorization).
		Get(ApiMe)

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
	reqBody := map[string]interface{}{
		"path":     path,
		"password": config.CONFIG.Password,
		"page":     1,
		"per_page": 1,
		"refresh":  false,
	}

	r, err := net.GlobalClient.R().
		SetHeader("Authorization", config.CONFIG.AListToken).
		SetBody(reqBody).
		Post(ApiFsGet)

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

	r, err := net.GlobalClient.R().
		SetHeaders(headers).
		SetFileReader("file", "name.png", f).
		Put(ApiFsForm)

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

	r, err := net.GlobalClient.R().
		SetHeaders(headers).
		SetBody(f).
		Put(ApiFsPut)

	if err != nil {
		return nil, err
	}
	if gjson.Get(string(r.Body()), "code").Int() != 200 {
		return nil, errors.New("alist /api/fs/put access error")
	}
	return r.Body(), nil
}

func FsSearch(name string) ([]byte, error) {
	reqBody := map[string]interface{}{
		"parent":   config.CONFIG.StoragePath,
		"password": config.CONFIG.Password,
		"page":     1,
		"per_page": 1,
		"scope":    2,
		"keywords": name,
	}

	r, err := net.GlobalClient.R().
		SetHeader("Authorization", config.CONFIG.AListToken).
		SetBody(reqBody).
		Post(ApiFsSearch)

	if err != nil {
		return nil, err
	}
	if gjson.Get(string(r.Body()), "code").Int() != 200 {
		return nil, errors.New(gjson.Get(string(r.Body()), "message").String())
	}
	if gjson.Get(string(r.Body()), "data.total").Int() == 0 {
		return nil, errors.New("alist no such item")
	}
	return r.Body(), nil
}

func Refresh() {
	reqBody := map[string]interface{}{
		"paths":     []string{config.CONFIG.StoragePath},
		"max_depth": 20,
	}

	_, _ = net.GlobalClient.R().
		SetHeader("Authorization", config.CONFIG.AListToken).
		SetBody(reqBody).
		Post(ApiIndexUpdate)
}
