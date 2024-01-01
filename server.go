package goeximbaselib

import (
	"github.com/exim-id/go-exim-base-lib/env"
	"github.com/exim-id/go-exim-base-lib/errors"
	"github.com/exim-id/go-exim-base-lib/mdw"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServerStart(endpoints func(*echo.Echo)) {
	e := echo.New()
	defer errors.ServerStop(e)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "[${time}] method=${method}, uri=${uri}, status=${status}\n"}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	e.HTTPErrorHandler = mdw.CustomHTTPErrorHandler
	endpoints(e)
	e.Logger.Fatal(e.Start(env.GetServerPort()))
}
