package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/sknutsen/harvestovertimeweb/lib"
	"github.com/sknutsen/harvestovertimeweb/models"
)

type Handler struct {
	Client       *http.Client
	ClientId     string
	ClientSecret string
	Port         string
	Database     string
	DbHost       string
	DbUser       string
	DbPass       string
	DbPort       string
	Secret       []byte
}

func (h *Handler) ConnectToCalendarDatabase() (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", h.DbHost, h.DbPort, h.DbUser, h.DbPass, h.Database)
	return sql.Open("postgres", connectionString)
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

func StoreSettingsAsCookie(c echo.Context, settings models.Settings) {
	data, err := json.Marshal(settings)
	if err != nil {
		println(err)
		return
	}

	cookieval := base64.StdEncoding.EncodeToString(data)

	SetCookie(c, "settings", cookieval, int(time.Hour*24*7))
}

func GetSettings(c echo.Context) models.Settings {
	cookie, err := c.Cookie("settings")
	if err != nil {
		fmt.Printf("failed getting settings cookie %s\n", err)
		return lib.DefaultSettings()
	}

	rawjson, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		fmt.Printf("Could not decode string %s\n", err)
	}

	var val models.Settings = models.Settings{}

	err = json.Unmarshal(rawjson, &val)
	if err != nil {
		fmt.Printf("failed parsing json %s\n", err)
		return lib.DefaultSettings()
	}

	return val
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
