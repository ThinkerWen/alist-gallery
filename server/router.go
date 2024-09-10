package server

import (
	"alist-gallery/server/handles"
	"github.com/labstack/echo/v4"
)

func RegisterFileSystem(e *echo.Echo) {
	g := e.Group("/fs")

	g.PUT("/form-gallery", handles.FormGallery)
	g.GET("/show-gallery/:name", handles.ShowImage)
}
