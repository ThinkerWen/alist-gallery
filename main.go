package main

import (
	"alist-gallery/config"
	"alist-gallery/server"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func registerRoutes(e *echo.Echo) {
	server.RegisterFileSystem(e)
}

func main() {
	e := echo.New()
	registerRoutes(e)
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("REQUEST",
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"status", c.Response().Status)
			return next(c)
		}
	})
	logger.Fatal(e.Start(fmt.Sprintf(":%d", config.CONFIG.Port)))
}
