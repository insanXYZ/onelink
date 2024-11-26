package controller

import (
	"radproject/model"
	"radproject/service"
	"radproject/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const url = "/user/site/"

type LinkController struct {
	linkService *service.LinkService
}

func NewLinkController(linkService *service.LinkService) *LinkController {
	return &LinkController{
		linkService: linkService,
	}
}

func (c *LinkController) CreateLink(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)
	req := new(model.CreateLinkRequest)
	err := ctx.Bind(req)
	urlPath := url + req.Site_Id
	if err != nil {
		return util.RedirectWithError(ctx, urlPath, err.Error())
	}
	err = c.linkService.CreateLink(ctx.Request().Context(), claims, req)
	if err != nil {
		return util.RedirectWithError(ctx, urlPath, err.Error())
	}
	return ctx.Redirect(303, urlPath)
}

func (c *LinkController) Delete(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)
	req := new(model.DeleteLinkRequest)
	err := ctx.Bind(req)
	if err != nil {
		return util.RedirectWithError(ctx, url+req.Site_Id, err.Error())
	}
	err = c.linkService.Delete(ctx.Request().Context(), claims, req)
	if err != nil {
		return util.RedirectWithError(ctx, url+req.Site_Id, err.Error())
	}
	return util.Redirect(ctx, url+req.Site_Id)
}

func (c *LinkController) Visit(ctx echo.Context) error {
	req := new(model.VisitLinkRequest)
	err := ctx.Bind(req)
	if err != nil {
		return err
	}
	link, err := c.linkService.Visit(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	return util.Redirect(ctx, link.Href)
}
