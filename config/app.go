package config

import (
	"database/sql"

	"radproject/controller"
	"radproject/middleware"
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
	LinkRepository := repository.NewLinkRepository()
	ClickRepository := repository.NewClikRepository()

	// serviceInit
	UserService := service.NewUserService(c.Validator, c.Db, UserRepository)
	LinkService := service.NewLinkService(c.Validator, c.Db, LinkRepository, SiteRepository)
	SiteService := service.NewSiteService(c.Validator, c.Db, SiteRepository, LinkRepository)
	clickService := service.NewClickService(c.Validator, c.Db, LinkRepository, SiteRepository, ClickRepository)

	// controllerInit
	UserController := controller.NewUserController(UserService)
	SiteController := controller.NewSiteController(SiteService)
	LinkController := controller.NewLinkController(LinkService)
	ClickController := controller.NewClickController(clickService)

	middleware := middleware.NewMiddleware(ClickController)

	routeConfig := &route.RouteConfig{
		Echo:           c.Echo,
		UserController: UserController,
		SiteController: SiteController,
		LinkController: LinkController,
		Middleware:     middleware,
	}

	route.InitRoute(routeConfig)
	c.Echo.Start(":1323")
}
