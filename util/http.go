package util

import (
	"radproject/model"

	"github.com/labstack/echo/v4"
)

func RedirectWithError(c echo.Context, statusCode int, url, errorMessage string) error {
	CreateFlashSession(c, model.SessionMessage, errorMessage, "error_message")
	return c.Redirect(statusCode, url)
}

func RenderViewHtml(c echo.Context, statusCode int, namefile string, data any) error {
	response := model.ResponseHttp{
		Data: data,
	}

	if flash, exist := GetFlashSession(c, model.SessionMessage, "error_message"); exist {
		response.Message = flash
	}

	return c.Render(statusCode, namefile, response)
}
