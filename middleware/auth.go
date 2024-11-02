package middleware

import (
	"errors"
	"fmt"
	"os"
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
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": claims["name"],
		"exp":  time.Now().Add(10 * time.Minute).Unix(),
	})

	return claim.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func isExpired(claims jwt.MapClaims) bool {
	return int64(claims["exp"].(float64)) <= time.Now().Unix()
}

func Guest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("go to middleware")
		fmt.Println("to " + c.Path())
		if token, exist := getCookiesSessionToken(c); exist {
			claims := jwt.MapClaims{}
			_, _ = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			// if err != nil {
			// 	if isExpired(claims) {
			// 		if token, err := refreshToken(claims); err == nil {
			// 			c.SetCookie(&http.Cookie{
			// 				Name:  model.SessionToken,
			// 				Value: token,
			// 				Path:  "/",
			// 			})
			// 		}
			// 	} else {
			// 		c.SetCookie(&http.Cookie{
			// 			Name:   model.SessionToken,
			// 			MaxAge: -1,
			// 		})
			// 		return next(c)
			// 	}
			// }
			return c.Redirect(303, "/")
		}

		return next(c)
	}
}
