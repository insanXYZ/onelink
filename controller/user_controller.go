package controller

import (
	"fmt"
	"net/http"

	"radproject/model"
	"radproject/service"
	"radproject/util"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

func (ctr *UserController) CreateLoginView(c echo.Context) error {
	return util.RenderViewHtml(c, 200, "login.html", nil)
}

func (ctr *UserController) CreateRegisterView(c echo.Context) error {
	return util.RenderViewHtml(c, 200, "register.html", nil)
}

func (ctr *UserController) CreateLandingPageView(c echo.Context) error {
  return util.RenderViewHtml(c, 200, "landing_page.html", nil)
}

func (ctr *UserController) Login(c echo.Context) error {
	req := new(model.LoginRequest)
	err := c.Bind(req)
	if err != nil {
		return util.RedirectWithError(c, 303, "/login", err.Error())
	}
	token, err := ctr.userService.Login(c.Request().Context(), req)
	if err != nil {
		return util.RedirectWithError(c, 303, "/login", err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = model.SessionToken
	cookie.Value = token
	cookie.Path = "/"

	c.SetCookie(cookie)

	return c.Redirect(303, "/dashboard")
}

func (ctr *UserController) Register(c echo.Context) error {
	req := new(model.RegisterRequest)
	err := c.Bind(req)
	if err != nil {
		return util.RedirectWithError(c, 303, "/register", err.Error())
	}
	err = ctr.userService.Register(c.Request().Context(), req)
	if err != nil {
		return util.RedirectWithError(c, 303, "/register", err.Error())
	}
	err = c.Redirect(303, "/login")
	if err != nil {
		fmt.Println("error redirect " + err.Error())
	}
	return err
}
