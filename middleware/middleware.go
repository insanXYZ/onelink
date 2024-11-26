package middleware

import "radproject/controller"

type Middleware struct {
	clickController *controller.ClickController
}

func NewMiddleware(clickController *controller.ClickController) *Middleware {
	return &Middleware{
		clickController: clickController,
	}
}
