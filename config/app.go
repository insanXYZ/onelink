package config

import (
	"database/sql"

	"radproject/controller"
	"radproject/repository"
	"radproject/route"
	"radproject/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type BootstrapConfigs struct {
	Validator *validator.Validate
	Echo      *echo.Echo
	Db        *sql.DB
}

func (c *BootstrapConfigs) Run() {
	// repositoryInit
	UserRepository := repository.NewUserRepository()

	// serviceInit
	UserService := service.NewUserService(c.Validator, c.Db, UserRepository)

	// controllerInit
	UserController := controller.NewUserController(UserService)

	routeConfig := &route.RouteConfig{
		Echo:           c.Echo,
		UserController: UserController,
	}

	route.InitRoute(routeConfig)
	c.Echo.Start(":1323")
}
