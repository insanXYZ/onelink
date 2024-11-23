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
	SiteRepository := repository.NewSiteRepository()

	// serviceInit
	UserService := service.NewUserService(c.Validator, c.Db, UserRepository)
	SiteService := service.NewSiteService(c.Validator, c.Db, SiteRepository)

	// controllerInit
	UserController := controller.NewUserController(UserService)
	SiteController := controller.NewSiteController(SiteService)

	routeConfig := &route.RouteConfig{
		Echo:           c.Echo,
		UserController: UserController,
		SiteController: SiteController,
	}

	route.InitRoute(routeConfig)
	c.Echo.Start(":1323")
}
