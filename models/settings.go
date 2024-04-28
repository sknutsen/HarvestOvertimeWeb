package models

import (
	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
)

type Settings struct {
	libmodels.Settings
	Years []int
}
