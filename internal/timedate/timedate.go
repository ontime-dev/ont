package timedate

import (
	"time"
)

func oneMonthFromNow(currentTime time.Time) time.Time {
	// Calculate the time for 1 month from now
	oneMonthLater := currentTime.AddDate(0, 1, 0)

	return oneMonthLater
}
