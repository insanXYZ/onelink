package controller

import (
	"errors"
	"net/http"

	"radproject/model"
	"radproject/model/message"
	"radproject/service"
	"radproject/util"

	"github.com/golang-jwt/jwt/v5"
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
	return util.RenderViewHtml(c, "login.html", nil)
}

func (ctr *UserController) CreateRegisterView(c echo.Context) error {
	return util.RenderViewHtml(c, "register.html", nil)
}

func (ctr *UserController) CreateLandingPageView(c echo.Context) error {

	auth := map[string]bool{
		"isLogin": false,
	}

	if _, err := c.Cookie(model.SessionToken); err == nil {
		auth["isLogin"] = true
	}

	return util.RenderViewHtml(c, "landing_page.html", auth)
}

func (ctr *UserController) CreateDashboardView(c echo.Context) error {
	return util.RenderViewHtml(c, "dashboard.html", nil)
}

func (ctr *UserController) CreateAccountView(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	resp, err := ctr.userService.GetAccount(c.Request().Context(), claims)
	if err != nil {
		return err
	}
	return util.RenderViewHtml(c, "account.html", *resp)
}

func (ctr *UserController) Login(c echo.Context) error {
	const message = "email or password is incorrect"

	req := new(model.LoginRequest)
	err := c.Bind(req)
	if err != nil {
		return util.RedirectWithError(c, "/login", message)
	}
	token, err := ctr.userService.Login(c.Request().Context(), req)
	if err != nil {
		return util.RedirectWithError(c, "/login", message)
	}

	cookie := new(http.Cookie)
	cookie.Name = model.SessionToken
	cookie.Value = token
	cookie.Path = "/"

	c.SetCookie(cookie)

	return util.Redirect(c, "/user/dashboard")
}

func (ctr *UserController) Register(c echo.Context) error {
	error := message.ERR_REGISTER

	req := new(model.RegisterRequest)
	err := c.Bind(req)
	if err != nil {
		return util.RedirectWithError(c, "/register", error.Error())
	}
	err = ctr.userService.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, message.ERR_REGISTER_EMAIL_USED) {
			error = message.ERR_REGISTER_EMAIL_USED
		}
		return util.RedirectWithError(c, "/register", error.Error())
	}
	return util.Redirect(c, "/login")
}

func (ctr *UserController) UpdateUser(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)

	req := new(model.UpdateUserRequest)
	err := c.Bind(req)
	if err != nil {
		return util.RedirectWithError(c, "/user/account", err.Error())
	}

	file, err := c.FormFile("image")
	if err == nil {
		req.Image = file
	}

	err = ctr.userService.UpdateUser(c.Request().Context(), claims, req)
	if err != nil {
		return util.RedirectWithError(c, "/user/account", err.Error())
	}

	return util.Redirect(c, "/user/account")
}

func (ctr *UserController) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:   model.SessionToken,
		MaxAge: -1,
		Path:   "/",
	}
	c.SetCookie(cookie)

	return util.Redirect(c, "/")
}
