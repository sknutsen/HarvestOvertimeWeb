package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/harvestovertimelib/v2"
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
	"github.com/sknutsen/harvestovertimeweb/routes"
	"github.com/sknutsen/harvestovertimeweb/view"
)

func (h *Handler) Index(c echo.Context) error {
	token, err := c.Cookie("accesstoken")

	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, routes.Auth)
	}

	settings := GetSettings(c)

	settings.AccessToken = token.Value

	userInfo, err := harvestovertimelib.GetUserInfo(h.Client, settings)

	if err != nil {
		userInfo = libmodels.UserInfo{}
	}

	settings.UserId = userInfo.ID

	tasks, err := harvestovertimelib.ListTasks(h.Client, settings)

	if err != nil {
		tasks = []libmodels.TaskDetails{}
	}

	component := view.Index(true, tasks, settings, userInfo)
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) Details(c echo.Context) error {
	token, err := c.Cookie("accesstoken")

	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	settings := GetSettings(c)

	settings.AccessToken = token.Value

	userInfo, err := harvestovertimelib.GetUserInfo(h.Client, settings)

	if err != nil {
		userInfo = libmodels.UserInfo{}
	}

	settings.UserId = userInfo.ID

	entries, err := harvestovertimelib.ListEntries(h.Client, settings)

	if err != nil {
		entries = libmodels.TimeEntries{}
	}

	filtered := harvestovertimelib.FilterTimeOffTasks(entries, settings)

	component := view.Details(filtered, userInfo)
	return component.Render(context.Background(), c.Response().Writer)
}
