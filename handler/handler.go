package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Client       *http.Client
	ClientId     string
	ClientSecret string
	Port         string
	Secret       []byte
}

func SetCookie(
	c echo.Context,
	key string,
	value string,
	maxAge int,
) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)
}

func GetSessionVar(c echo.Context, key string) interface{} {
	sess, _ := session.Get("session", c)

	return sess.Values[key]
}

func SetSessionVar(c echo.Context, key string, value string, maxAge int) {
	sess, _ := session.Get("session", c)

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
	}
}
