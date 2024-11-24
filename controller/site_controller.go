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
}

func NewSiteController(siteService *service.SiteService) *SiteController {
	return &SiteController{
		siteService: siteService,
	}
}

func (c *SiteController) CreateSiteView(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)

	sites, err := c.siteService.GetAllSites(ctx.Request().Context(), claims)
	if err != nil {
		log.Println(err.Error())
	}

	return util.RenderViewHtml(ctx, "site.html", sites)
}

func (c *SiteController) CreateEditSiteView(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)
	id := ctx.Param("id")

	site, err := c.siteService.GetSite(ctx.Request().Context(), claims, id)
	if err != nil {
		return util.RedirectWithError(ctx, "/user/site", err.Error())
	}

	return util.RenderViewHtml(ctx, "edit_site.html", *site)
}

func (c *SiteController) CreateSite(ctx echo.Context) error {
	const message = "create site failed"

	claims := ctx.Get("user").(jwt.MapClaims)
	req := new(model.CreateSiteRequest)
	err := ctx.Bind(req)
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, "/user/site", message)
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, "/user/site", message)
	}
	req.Image = file

	err = c.siteService.CreateSite(ctx.Request().Context(), claims, req)
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, "/user/site", message)
	}

	return util.Redirect(ctx, "/user/site")
}

func (c *SiteController) Delete(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)

	req := new(model.DeleteSiteRequest)
	err := ctx.Bind(req)
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, "/user/site", err.Error())
	}

	err = c.siteService.DeleteSite(ctx.Request().Context(), claims, req)
	if err != nil {
		log.Println(err.Error())
		return util.RedirectWithError(ctx, "/user/site", err.Error())
	}

	return util.Redirect(ctx, "/user/site")
}

func (c *SiteController) Update(ctx echo.Context) error {
	claims := ctx.Get("user").(jwt.MapClaims)
	req := new(model.UpdateSiteRequest)
	err := ctx.Bind(req)
	if err != nil {
		return util.RedirectWithError(ctx, "/user/site", err.Error())
	}
	file, err := ctx.FormFile("image")
	if err == nil {
		req.Image = file
	}
	err = c.siteService.UpdateSite(ctx.Request().Context(), claims, req)
	if err != nil {
		return util.RedirectWithError(ctx, "/user/site", err.Error())
	}
	return util.Redirect(ctx, "/user/site/"+req.Id)
}
