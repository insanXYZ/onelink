package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"radproject/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func getCookiesSessionToken(c echo.Context) (string, bool) {
	cookies := c.Cookies()

	var token string

	for _, c := range cookies {
		if c.Name == model.SessionToken {
			token = c.Value
			return token, true
		}
	}

	return token, false
}

func refreshToken(claims jwt.MapClaims) (string, error) {

	exp, _ := strconv.Atoi(os.Getenv("JWT_EXP"))

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": claims["name"],
		"id":   claims["id"],
		"exp":  time.Now().Add(time.Duration(exp) * time.Minute).Unix(),
	})

	return claim.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func setContextUser(c echo.Context, claims jwt.MapClaims) {
	c.Set("user", claims)
}

func isExpired(claims jwt.MapClaims) bool {
	return int64(claims["exp"].(float64)) <= time.Now().Unix()
}

func (m *Middleware) Guest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if token, exist := getCookiesSessionToken(c); exist {
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil {
				if isExpired(claims) {
					if rtoken, err := refreshToken(claims); err == nil {
						c.SetCookie(&http.Cookie{
							Name:  model.SessionToken,
							Value: rtoken,
							Path:  "/",
						})
						return c.Redirect(303, "/")
					}
				}
				c.SetCookie(&http.Cookie{
					Name:   model.SessionToken,
					MaxAge: -1,
					Path:   "/",
				})
				return next(c)
			}

			return c.Redirect(303, "/")

		}

		return next(c)
	}
}

func (m *Middleware) IsLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if token, exist := getCookiesSessionToken(c); exist {
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil {
				if isExpired(claims) {
					if rtoken, err := refreshToken(claims); err == nil {
						cookie := &http.Cookie{
							Name:  model.SessionToken,
							Value: rtoken,
							Path:  "/",
						}
						c.SetCookie(cookie)
						setContextUser(c, claims)
						return next(c)
					}
				}

				c.SetCookie(&http.Cookie{
					Name:   model.SessionToken,
					MaxAge: -1,
					Path:   "/",
				})
				return c.Redirect(303, "/login")
			}
			setContextUser(c, claims)
			return next(c)
		}
		return c.Redirect(303, "/login")
	}
}
