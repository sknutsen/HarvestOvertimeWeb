package lib

import (
	"fmt"
	"time"

	"github.com/sknutsen/harvestovertimelib/v2/lib"
	"github.com/sknutsen/harvestovertimelib/v2/models"
)

func DateToString(date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()

	var dayString string
	var monthString string

	if day < 10 {
		dayString = fmt.Sprintf("0%d", day)
	} else {
		dayString = fmt.Sprint(day)
	}

	if month < 10 {
		monthString = fmt.Sprintf("0%d", month)
	} else {
		monthString = fmt.Sprintf("%d", month)
	}

	dateAsString := fmt.Sprintf("%d-%s-%s", year, monthString, dayString)

	return dateAsString
}

func TaskIsSelected(task models.TaskDetails, tasks []models.Task) bool {
	for _, t := range tasks {
		if t.ID == task.Task.ID {
			return true
		}
	}

	return false
}

func SumHoursFromEntries(entries []models.TimeEntry) float64 {
	var sum float64 = 0

	for _, e := range entries {
		sum += e.Hours
	}

	return sum
}

func DefaultSettings() models.Settings {
	return models.Settings{
		FromDate:                 fmt.Sprintf("%d-01-01", time.Now().Year()-2),
		ToDate:                   DateToString(time.Now()),
		DaysInWeek:               5,
		WorkDayHours:             7.5,
		SimulateFullWeekAtToDate: true,
		TimeOffTasks: []models.Task{
			{
				ID: 10882012,
			},
		},
		WorkDays: []time.Weekday{
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
		},
	}
}

func DateIsValid(settings models.Settings, date string) bool {
	if DateIsInSpan(settings, date) {
		return DateIsWorkday(settings, date)
	}

	return false
}

func DateIsInSpan(settings models.Settings, date string) bool {
	from := lib.ParseDateString(settings.FromDate)
	to := lib.ParseDateString(settings.ToDate)
	datetime := lib.ParseDateString(date)

	return datetime.After(from) && datetime.Before(to)
}

func DateIsWorkday(settings models.Settings, date string) bool {
	datetime := lib.ParseDateString(date)

	for _, d := range settings.WorkDays {
		if datetime.Weekday() == d {
			return true
		}
	}

	return false
}
