package util

import (
	"radproject/model"
	"time"

	"github.com/labstack/echo/v4"
)

func Redirect(c echo.Context, url string) error {
	c.Response().Header().Add("Date", time.Now().Format(time.DateOnly))
	return c.Redirect(303, url)
}

func RedirectWithError(c echo.Context, url, errorMessage string) error {
	CreateFlashSession(c, model.SessionMessage, errorMessage, "error_message")
	return Redirect(c, url)
}

func RenderViewHtml(c echo.Context, namefile string, data any) error {
	response := model.ResponseHttp{
		Data: data,
	}

	if flash, exist := GetFlashSession(c, model.SessionMessage, "error_message"); exist {
		response.Message = flash
	}

	return c.Render(200, namefile, response)
}
