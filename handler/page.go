package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	fromDate := settings.FromDate

	settings.FromDate = fmt.Sprintf("%d-01-01", time.Now().Year()-2)
	settings.ToDate = fmt.Sprintf("%d-12-31", time.Now().Year())

	tasks, err := harvestovertimelib.ListTasks(h.Client, settings)

	if err != nil {
		tasks = []libmodels.TaskDetails{}
	}

	settings.FromDate = fromDate

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

	holidays, _ := h.GetCalendarEvents()

	entries.TimeEntries = append(entries.TimeEntries, ConvertHolidaysToTimeEntries(settings, holidays)...)

	filtered := harvestovertimelib.FilterTimeOffTasks(entries, settings)

	component := view.Details(filtered, userInfo)
	return component.Render(context.Background(), c.Response().Writer)
}
