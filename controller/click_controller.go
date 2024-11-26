package controller

import (
	"context"
	"radproject/model"
	"radproject/service"

	"github.com/labstack/echo/v4"
)

type ClickController struct {
	clickService *service.ClickService
}

func NewClickController(clickService *service.ClickService) *ClickController {
	return &ClickController{
		clickService: clickService,
	}
}

func (c *ClickController) Visit(ctx echo.Context) error {
	req := new(model.VisitDestination)
	err := ctx.Bind(req)
	if err != nil {
		return err
	}

	if req.Domain == "favicon.ico" {
		return nil
	}

	err = c.clickService.Visit(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}
