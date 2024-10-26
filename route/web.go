package route

import (
	"radproject/controller"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	Echo           *echo.Echo
	UserController *controller.UserController
}

func InitRoute(config *RouteConfig) {
	e := config.Echo
	e.Static("/storage", "storage")
	e.POST("/login", config.UserController.Login)
	e.POST("/register", config.UserController.Register)
	e.GET("/login", config.UserController.CreateLoginView)
	e.GET("/register", config.UserController.CreateRegisterView)
}
