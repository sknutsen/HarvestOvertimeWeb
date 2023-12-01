package models

type GetHoursRequest struct {
	TimeOffTasks []int `form:"timeOffTasks"`
}
