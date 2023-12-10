package models

import (
	"time"

	libmodels "github.com/sknutsen/harvestovertimelib/v2/models"
)

type TimeEntriesByDate struct {
	Date        time.Time
	TimeEntries []libmodels.TimeEntry
}
