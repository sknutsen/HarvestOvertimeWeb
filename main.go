package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sknutsen/harvestovertimelib"
	libmodels "github.com/sknutsen/harvestovertimelib/models"
	"github.com/sknutsen/harvestovertimeweb/models"
	"github.com/sknutsen/harvestovertimeweb/view"
)

var client *http.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Could not load .env file. Err: %s\n", err)
	}

	clientId := os.Getenv("HARVEST_CLIENT_ID")
	clientSecret := os.Getenv("HARVEST_CLIENT_SECRET")
	port := os.Getenv("PORT")

	client = &http.Client{Timeout: 10 * time.Second}

	e := echo.New()

	e.Use(middleware.Logger())

	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		_, err := c.Cookie("accesstoken")

		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/hours")
		}

		settings := libmodels.Settings{}

		component := view.Index(true, settings)
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/hours", func(c echo.Context) error {
		url := fmt.Sprintf("https://id.getharvest.com/oauth2/authorize?client_id=%s&response_type=code", clientId)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	e.POST("/hours/get", func(c echo.Context) error {
		refreshTokenCookie, _ := c.Cookie("refreshtoken")

		// c.FormValue("")
		token, err := refreshToken(client, refreshTokenCookie.Value, clientId, clientSecret)
		if err != nil {
			component := view.Index(false, libmodels.Settings{})
			return component.Render(context.Background(), c.Response().Writer)

		}

		SetCookie(c, "accesstoken", token.AccessToken, int(token.ExpiresIn))
		SetCookie(c, "refreshtoken", token.RefreshToken, int(token.ExpiresIn))

		settings := libmodels.Settings{
			AccessToken:     token.AccessToken,
			CarryOverTime:   0,
			OnlyCurrentYear: true,
			TimeOffTasks: []libmodels.Task{
				{
					Id:   10882012,
					Name: "Avspasering",
				},
			},
		}

		entries, _ := harvestovertimelib.ListEntries(client, settings)
		hours := harvestovertimelib.GetTotalOvertime(entries, settings)

		component := view.Hours(hours)
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/signin_callback", func(c echo.Context) error {
		authCode := c.QueryParam("code")

		token, err := newToken(client, authCode, clientId, clientSecret)
		if err != nil {
			component := view.Index(false, libmodels.Settings{})
			return component.Render(context.Background(), c.Response().Writer)

		}

		var buf bytes.Buffer

		err = gob.NewEncoder(&buf).Encode(&token)
		if err != nil {
			return err
		}

		SetCookie(c, "accesstoken", token.AccessToken, int(token.ExpiresIn))
		SetCookie(c, "refreshtoken", token.RefreshToken, int(token.ExpiresIn))

		fmt.Printf("Token expires in: %d\n", token.ExpiresIn)

		return c.Redirect(http.StatusPermanentRedirect, "/")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func newToken(
	client *http.Client,
	code string,
	clientId string,
	clientSecret string,
) (models.HarvestToken, error) {
	body := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	url := "https://id.getharvest.com/api/v2/oauth2/token"

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body.Encode()))
	if err != nil {
		println("Error creating request: " + err.Error())
		return models.HarvestToken{}, err
	}

	return getToken(req, client)
}

func refreshToken(
	client *http.Client,
	refreshToken string,
	clientId string,
	clientSecret string,
) (models.HarvestToken, error) {
	body := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	url := "https://id.getharvest.com/api/v2/oauth2/token"

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body.Encode()))
	if err != nil {
		println("Error creating request: " + err.Error())
		return models.HarvestToken{}, err
	}

	return getToken(req, client)
}

func getToken(req *http.Request, client *http.Client) (models.HarvestToken, error) {
	var token models.HarvestToken

	req.Header.Add("User-Agent", os.Getenv("USER_AGENT"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		println("Error sending request: " + err.Error())
		return models.HarvestToken{}, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		println("Error decoding response: " + err.Error())
		return models.HarvestToken{}, err
	}

	return token, nil
}

func SetCookie(
	c echo.Context,
	key string,
	value string,
	maxAge int,
) {
	cookie := &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: maxAge,
		Path:   "/",
	}

	c.SetCookie(cookie)
}
