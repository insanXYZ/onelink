package route

import (
	"radproject/controller"
	"radproject/middleware"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	Echo           *echo.Echo
	UserController *controller.UserController
	SiteController *controller.SiteController
	LinkController *controller.LinkController
	Middleware     *middleware.Middleware
}

func InitRoute(config *RouteConfig) {
	e := config.Echo

	e.Static("/storage", "storage")

	e.GET("/", config.UserController.CreateLandingPageView)

	guest := e.Group("", config.Middleware.Guest)
	guest.POST("/login", config.UserController.Login)
	guest.POST("/register", config.UserController.Register)
	guest.GET("/login", config.UserController.CreateLoginView)
	guest.GET("/register", config.UserController.CreateRegisterView)

	user := e.Group("/user", config.Middleware.IsLogin)
	user.GET("/dashboard", config.UserController.CreateDashboardView)
	user.GET("/account", config.UserController.CreateAccountView)
	user.PATCH("/account", config.UserController.UpdateUser)
	user.GET("/logout", config.UserController.Logout)

	site := user.Group("/site")
	site.GET("/", func(c echo.Context) error {
		return c.Redirect(303, "/user/site")
	})
	site.GET("", config.SiteController.CreateSiteView)
	site.POST("", config.SiteController.CreateSite)
	site.DELETE("/:id", config.SiteController.Delete)
	site.GET("/:id", config.SiteController.CreateEditSiteView)
	site.PATCH("/:id", config.SiteController.Update)

	site.POST("/:id", config.LinkController.CreateLink)
	site.DELETE("/:site_id/:link_id", config.LinkController.Delete)

	click := e.Group("", config.Middleware.Click)
	click.GET("/:domain_site", config.SiteController.CreatePublishSiteView)
	click.GET("/link/visit/:id", config.LinkController.Visit)
}
