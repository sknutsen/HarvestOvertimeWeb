package models

import (
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
)

type ClientState struct {
	Tasks    []libmodels.TaskDetails
	UserInfo libmodels.UserInfo
	Years    []int
}
