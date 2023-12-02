package models

type GetHoursRequest struct {
	TimeOffTasks             []int   `form:"timeOffTasks"`
	FromDate                 string  `form:"fromDate"`
	ToDate                   string  `form:"toDate"`
	DaysInWeek               int     `form:"daysInWeek"`
	WorkDayHours             float32 `form:"workDayHours"`
	SimulateFullWeekAtToDate bool    `form:"simulateFullWeekAtToDate"`
}
