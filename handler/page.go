package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/harvestovertimelib/v2"
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
	"github.com/sknutsen/harvestovertimeweb/lib"
	"github.com/sknutsen/harvestovertimeweb/routes"
	"github.com/sknutsen/harvestovertimeweb/view"
)

func (h *Handler) Index(c echo.Context) error {
	token, err := c.Cookie("accesstoken")

	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, routes.Auth)
	}

	settings := libmodels.Settings{
		AccessToken:              token.Value,
		FromDate:                 fmt.Sprintf("%d-01-01", time.Now().Year()),
		ToDate:                   lib.DateToString(time.Now()),
		DaysInWeek:               5,
		WorkDayHours:             7.5,
		SimulateFullWeekAtToDate: true,
		TimeOffTasks: []libmodels.Task{
			{
				ID: 10882012,
			},
		},
	}

	userInfo, err := harvestovertimelib.GetUserInfo(h.Client, settings)

	if err != nil {
		userInfo = libmodels.UserInfo{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	settings.UserId = userInfo.ID

	tasks, err := harvestovertimelib.ListTasks(h.Client, settings)

	if err != nil {
		tasks = []libmodels.TaskDetails{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	component := view.Index(true, tasks, settings, userInfo)
	return component.Render(context.Background(), c.Response().Writer)
}

func (h *Handler) Details(c echo.Context) error {
	token, err := c.Cookie("accesstoken")

	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	settings := libmodels.Settings{
		AccessToken:              token.Value,
		FromDate:                 fmt.Sprintf("%d-01-01", time.Now().Year()),
		ToDate:                   fmt.Sprintf("%d-12-31", time.Now().Year()),
		DaysInWeek:               5,
		WorkDayHours:             7.5,
		SimulateFullWeekAtToDate: true,
		TimeOffTasks: []libmodels.Task{
			{
				ID: 10882012,
			},
		},
	}

	userInfo, err := harvestovertimelib.GetUserInfo(h.Client, settings)

	if err != nil {
		userInfo = libmodels.UserInfo{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	settings.UserId = userInfo.ID

	entries, err := harvestovertimelib.ListEntries(h.Client, settings)

	if err != nil {
		entries = libmodels.TimeEntries{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	filtered := harvestovertimelib.FilterTimeOffTasks(entries, settings)

	component := view.Details(filtered, userInfo)
	return component.Render(context.Background(), c.Response().Writer)
}
