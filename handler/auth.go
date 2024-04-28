package handler

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

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/harvestovertimeweb/models"
	"github.com/sknutsen/harvestovertimeweb/routes"
	"github.com/sknutsen/harvestovertimeweb/view"
)

func (h *Handler) Auth(c echo.Context) error {
	url := fmt.Sprintf("https://id.getharvest.com/oauth2/authorize?client_id=%s&response_type=code", h.ClientId)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) Callback(c echo.Context) error {
	authCode := c.QueryParam("code")

	token, err := newToken(h.Client, authCode, h.ClientId, h.ClientSecret)
	if err != nil {
		component := view.Index(false, models.ClientState{}, models.Settings{})
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

	return c.Redirect(http.StatusPermanentRedirect, routes.Index)
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
