package lib

import (
	"fmt"
	"time"

	"github.com/sknutsen/harvestovertimelib/v2/lib"
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
	"github.com/sknutsen/harvestovertimeweb/models"
	"github.com/sknutsen/harvestovertimeweb/tasks"
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

func TaskIsSelected(task libmodels.TaskDetails, taskList []libmodels.Task) bool {
	for _, t := range taskList {
		if t.ID == task.Task.ID {
			return true
		}
	}

	return false
}

func SumHoursFromEntries(entries []libmodels.TimeEntry) float64 {
	var sum float64 = 0

	for _, e := range entries {
		if e.Task.ID != tasks.HolidayTask {
			sum += e.Hours
		}
	}

	return sum
}

func DefaultSettings() models.Settings {
	return models.Settings{
		Settings: libmodels.Settings{
			FromDate:                 "2021-01-01",
			ToDate:                   DateToString(time.Now()),
			DaysInWeek:               5,
			WorkDayHours:             7.5,
			SimulateFullWeekAtToDate: true,
			TimeOffTasks: []libmodels.Task{
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
		},
		Years: []int{
			2024,
			2025,
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

	onOrAfterFrom := datetime.Compare(from) > -1
	onOrBeforeTo := datetime.Compare(to) < 1

	return onOrAfterFrom && onOrBeforeTo
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
