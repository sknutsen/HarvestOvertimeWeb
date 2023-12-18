package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/harvestovertimelib/v2"
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
	"github.com/sknutsen/harvestovertimeweb/models"
	"github.com/sknutsen/harvestovertimeweb/view"
)

func (h *Handler) GetOvertimeHours(c echo.Context) error {
	refreshTokenCookie, _ := c.Cookie("refreshtoken")

	var getHoursRequest models.GetHoursRequest

	err := c.Bind(&getHoursRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	fmt.Printf("Time off tasks len: %d\n", len(getHoursRequest.TimeOffTasks))
	token, err := refreshToken(h.Client, refreshTokenCookie.Value, h.ClientId, h.ClientSecret)
	if err != nil {
		component := view.Index(false, []libmodels.TaskDetails{}, libmodels.Settings{}, libmodels.UserInfo{})
		return component.Render(context.Background(), c.Response().Writer)

	}

	SetCookie(c, "accesstoken", token.AccessToken, int(token.ExpiresIn))
	SetCookie(c, "refreshtoken", token.RefreshToken, int(token.ExpiresIn))

	settings := libmodels.Settings{
		AccessToken:              token.AccessToken,
		CarryOverTime:            0,
		WorkDayHours:             getHoursRequest.WorkDayHours,
		DaysInWeek:               getHoursRequest.DaysInWeek,
		FromDate:                 getHoursRequest.FromDate,
		ToDate:                   getHoursRequest.ToDate,
		SimulateFullWeekAtToDate: getHoursRequest.SimulateFullWeekAtToDate,
		WorkDays: []time.Weekday{
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
		},
		TimeOffTasks: []libmodels.Task{},
	}

	fmt.Printf("Current date: %s\n", settings.ToDate)

	for _, taskId := range getHoursRequest.TimeOffTasks {
		settings.TimeOffTasks = append(settings.TimeOffTasks, libmodels.Task{
			ID: uint64(taskId),
		})
	}

	userInfo, err := harvestovertimelib.GetUserInfo(h.Client, settings)

	if err != nil {
		userInfo = libmodels.UserInfo{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	settings.UserId = userInfo.ID

	entries, _ := harvestovertimelib.ListEntries(h.Client, settings)
	hours := harvestovertimelib.GetTotalOvertime(entries, settings)

	component := view.Hours(hours)
	return component.Render(context.Background(), c.Response().Writer)
}
