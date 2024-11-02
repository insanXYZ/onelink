package route

import (
	"radproject/controller"
	"radproject/middleware"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	Echo           *echo.Echo
	UserController *controller.UserController
}

func InitRoute(config *RouteConfig) {
	e := config.Echo
	e.Static("/storage", "storage")

	e.GET("/", config.UserController.CreateLandingPageView)

	guest := e.Group("", middleware.Guest)
	guest.POST("/login", config.UserController.Login)
	guest.POST("/register", config.UserController.Register)
	guest.GET("/login", config.UserController.CreateLoginView)
	guest.GET("/register", config.UserController.CreateRegisterView)
}
