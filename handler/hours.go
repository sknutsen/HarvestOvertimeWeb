package handler

import (
	"context"
	"fmt"
	"net/http"

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
		component := view.Index(false, models.ClientState{}, models.Settings{})
		return component.Render(context.Background(), c.Response().Writer)

	}

	SetCookie(c, "accesstoken", token.AccessToken, int(token.ExpiresIn))
	SetCookie(c, "refreshtoken", token.RefreshToken, int(token.ExpiresIn))

	settings := models.Settings{
		Settings: libmodels.Settings{
			AccessToken:              token.AccessToken,
			CarryOverTime:            getHoursRequest.CarryOverTime,
			WorkDayHours:             getHoursRequest.WorkDayHours,
			DaysInWeek:               len(getHoursRequest.Workdays),
			FromDate:                 getHoursRequest.FromDate,
			ToDate:                   getHoursRequest.ToDate,
			SimulateFullWeekAtToDate: getHoursRequest.SimulateFullWeekAtToDate,
			WorkDays:                 getHoursRequest.Workdays,
			TimeOffTasks:             []libmodels.Task{},
		},
		Years: getHoursRequest.Years,
	}

	for _, taskId := range getHoursRequest.TimeOffTasks {
		settings.TimeOffTasks = append(settings.TimeOffTasks, libmodels.Task{
			ID: uint64(taskId),
		})
	}

	userInfo, err := harvestovertimelib.GetUserInfo(h.Client, settings.Settings)

	if err != nil {
		userInfo = libmodels.UserInfo{}
		// return c.Redirect(http.StatusTemporaryRedirect, "/hours")
	}

	settings.UserId = userInfo.ID

	StoreSettingsAsCookie(c, settings)

	entries, err := harvestovertimelib.ListEntries(h.Client, settings.Settings)

	if err != nil {
		entries = libmodels.TimeEntries{}
	}

	holidays, err := h.GetCalendarEvents()

	if err == nil {
		entries.TimeEntries = append(entries.TimeEntries, ConvertHolidaysToTimeEntries(settings, holidays)...)
	}

	hours := harvestovertimelib.GetTotalOvertime(entries, settings.Settings)

	component := view.Hours(hours)
	return component.Render(context.Background(), c.Response().Writer)
}
