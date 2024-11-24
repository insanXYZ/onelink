package controller

import (
	"log"
	"radproject/model"
	"radproject/service"
	"radproject/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type SiteController struct {
	siteService *service.SiteService
	linkService *service.LinkService
}

func NewSiteController(siteService *service.SiteService, linkService *service.LinkService) *SiteController {
	return &SiteController{
		siteService: siteService,
		linkService: linkService,
	}
}

func (c *SiteController) CreateSiteView(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)

	sites, err := c.siteService.GetAllSites(ctx.Request().Context(), claims)
	if err != nil {
		log.Println(err.Error())
	}

	return util.RenderViewHtml(ctx, 200, "site.html", sites)
}

func (c *SiteController) CreateEditSiteView(ctx echo.Context) error {
	// claims := ctx.Get("user").(jwt.MapClaims)
	id := ctx.Param("id")
	return util.RenderViewHtml(ctx, 200, "edit_site.html", id)
}

func (c *SiteController) CreateSite(ctx echo.Context) error {
	const message = "create site failed"

	claims := ctx.Get("user").(jwt.MapClaims)
	req := new(model.CreateSiteRequest)
	err := ctx.Bind(req)
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, 200, "/user/site", message)
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, 200, "/user/site", message)
	}
	req.Image = file

	err = c.siteService.CreateSite(ctx.Request().Context(), claims, req)
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, 200, "/user/site", message)
	}

	return ctx.Redirect(303, "/user/site")
}

func (c *SiteController) Delete(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)

	req := new(model.DeleteSiteRequest)
	err := ctx.Bind(req)
	if err != nil {
		return util.RedirectWithError(ctx, 200, "/user/site", err.Error())
	}

	err = c.siteService.DeleteSite(ctx.Request().Context(), claims, req)
	if err != nil {
		return util.RedirectWithError(ctx, 200, "/user/site", err.Error())
	}

	return ctx.Redirect(303, "/user/site")
}
