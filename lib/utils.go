package lib

import (
	"fmt"
	"time"

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

func TaskIsSelected(task models.Task, tasks []models.Task) bool {
	for _, t := range tasks {
		if t.ID == task.ID {
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
