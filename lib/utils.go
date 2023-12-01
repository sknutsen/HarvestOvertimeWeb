package lib

import (
	"fmt"
	"time"
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
