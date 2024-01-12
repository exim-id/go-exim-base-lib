package mdw

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Info("Code", code)
	c.Logger().Info(err)
}

func ErrorRecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Error", r)
				log.Panicln("Stack trace", string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
			}
		}()
		if err := next(c); err != nil {
			log.Printf("Error occurred: %v", err)
			he, _ := err.(*echo.HTTPError)
			return c.JSON(he.Code, map[string]interface{}{"message": he.Message})
		}
		return nil
	}
}
