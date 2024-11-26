package middleware

import (
	"github.com/labstack/echo/v4"
)

func (m *Middleware) Click(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		go m.clickController.Visit(c)

		return next(c)
	}
}
