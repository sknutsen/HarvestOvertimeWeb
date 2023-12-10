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

	tasks, err := harvestovertimelib.ListTasks(h.Client, settings)

	if err != nil {
		tasks = []libmodels.Task{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	component := view.Index(true, tasks, settings)
	return component.Render(context.Background(), c.Response().Writer)
}
