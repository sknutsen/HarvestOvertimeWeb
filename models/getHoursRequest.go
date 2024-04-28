package models

import "time"

type GetHoursRequest struct {
	TimeOffTasks             []int          `form:"timeOffTasks"`
	Years                    []int          `form:"years"`
	FromDate                 string         `form:"fromDate"`
	ToDate                   string         `form:"toDate"`
	Workdays                 []time.Weekday `form:"workdays"`
	DaysInWeek               int            `form:"daysInWeek"`
	WorkDayHours             float32        `form:"workDayHours"`
	SimulateFullWeekAtToDate bool           `form:"simulateFullWeekAtToDate"`
}
