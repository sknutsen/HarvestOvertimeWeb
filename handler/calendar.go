package handler

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sknutsen/harvestovertimelib/v2/models"
	"github.com/sknutsen/harvestovertimeweb/lib"
)

type Holiday struct {
	Id          int
	Calendar    string
	Date        string
	Description string
}

func (h *Handler) GetCalendarEvents() ([]Holiday, error) {
	var holidays []Holiday

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", h.DbHost, h.DbPort, h.DbUser, h.DbPass, h.Database)

	conn, err := sql.Open("postgres", connectionString)

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

func ConvertHolidaysToTimeEntries(settings models.Settings, holidays []Holiday) []models.TimeEntry {
	var entries []models.TimeEntry = []models.TimeEntry{}

	for _, holiday := range holidays {
		if lib.DateIsValid(settings, holiday.Date) {
			entries = append(entries, models.TimeEntry{
				SpentDate: holiday.Date,
				Hours:     float64(settings.WorkDayHours),
				Task: models.Task{
					ID:   1,
					Name: fmt.Sprintf("%s - %s", holiday.Calendar, holiday.Description),
				},
			})
		}
	}

	return entries
}
