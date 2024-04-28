package handler

import (
	"fmt"
	"slices"
	"strconv"

	_ "github.com/lib/pq"
	liblib "github.com/sknutsen/harvestovertimelib/v2/lib"
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
	"github.com/sknutsen/harvestovertimeweb/lib"
	"github.com/sknutsen/harvestovertimeweb/models"
)

type Holiday struct {
	Id          int
	Calendar    string
	Date        string
	Description string
}

func (h *Handler) GetCalendarEvents() ([]Holiday, error) {
	var holidays []Holiday

	conn, err := h.ConnectToCalendarDatabase()

	if err != nil {
		return holidays, err
	}

	defer conn.Close()

	rows, err := conn.Query("SELECT calendar, date, description FROM dates WHERE calendar = 'NO'")

	if err != nil {
		return holidays, err
	}

	defer rows.Close()

	holidays = []Holiday{}

	for rows.Next() {
		var calendar string
		var date string
		var description string

		err = rows.Scan(&calendar, &date, &description)

		if err != nil {
			return holidays, err
		}

		holidays = append(holidays, Holiday{
			Calendar:    calendar,
			Date:        date,
			Description: description,
		})
	}

	return holidays, nil
}

func (h *Handler) GetCalendarYears() ([]int, error) {
	var years []int

	conn, err := h.ConnectToCalendarDatabase()

	if err != nil {
		return years, err
	}

	defer conn.Close()

	rows, err := conn.Query("SELECT date FROM dates WHERE calendar = 'NO' ORDER BY date ASC")

	if err != nil {
		return years, err
	}

	defer rows.Close()

	years = []int{}

	for rows.Next() {
		var date string

		err = rows.Scan(&date)

		if err != nil {
			return years, err
		}

		year, err := strconv.Atoi(date[:4])

		if err != nil {
			return years, err
		}

		if !slices.Contains(years, year) {
			years = append(years, year)
		}
	}

	return years, nil
}

func ConvertHolidaysToTimeEntries(settings models.Settings, holidays []Holiday) []libmodels.TimeEntry {
	var entries []libmodels.TimeEntry = []libmodels.TimeEntry{}

	for _, holiday := range holidays {
		if slices.Contains(settings.Years, liblib.ParseDateString(holiday.Date).Year()) && lib.DateIsValid(settings, holiday.Date) {
			entries = append(entries, libmodels.TimeEntry{
				SpentDate: holiday.Date,
				Hours:     float64(settings.WorkDayHours),
				Task: libmodels.Task{
					ID:   1,
					Name: fmt.Sprintf("%s - %s", holiday.Calendar, holiday.Description),
				},
			})
		}
	}

	return entries
}
