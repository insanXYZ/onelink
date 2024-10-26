package util

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var Session = sessions.NewCookieStore([]byte("yangtautauaja"))

func CreateSession(c echo.Context, name string, values map[any]any) error {
	s, err := Session.Get(c.Request(), name)
	if err != nil {
		return err
	}
	s.Values = values
	return s.Save(c.Request(), c.Response())
}

func CreateFlashSession(c echo.Context, name string, value any, varflash ...string) error {
	s, err := Session.Get(c.Request(), name)
	if err != nil {
		return err
	}
	s.AddFlash(value, varflash[0])
	return s.Save(c.Request(), c.Response())
}

func GetFlashSession(c echo.Context, name, varflash string) (string, bool) {
	s, err := Session.Get(c.Request(), name)
	if err != nil {
		panic(err.Error())
	}
	f := s.Flashes(varflash)
	s.Save(c.Request(), c.Response())
	if len(f) > 0 {
		return f[0].(string), true
	}
	return "", false
}
