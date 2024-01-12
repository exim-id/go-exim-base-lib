package goeximbaselib

import (
	"github.com/exim-id/go-exim-base-lib/env"
	"github.com/exim-id/go-exim-base-lib/errors"
	"github.com/exim-id/go-exim-base-lib/mdw"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ServerStart(endpoints func(*echo.Echo)) {
	e := echo.New()
	defer errors.ServerStop(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "[${time}] method=${method}, uri=${uri}, status=${status}\n"}))
	e.Use(mdw.ErrorRecoverMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.DELETE, echo.GET, echo.HEAD, echo.PATCH, echo.POST, echo.PUT},
	}))
	e.HTTPErrorHandler = mdw.CustomHTTPErrorHandler
	endpoints(e)
	e.Logger.Fatal(e.Start(env.GetServerPort()))
}
