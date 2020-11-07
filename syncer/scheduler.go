package syncer

import (
	"errors"

	"github.com/go-co-op/gocron"
)

// SetScheduler translate the syncersaccess
func (a *Access) SetScheduler(scheduler *gocron.Scheduler) (err error) {
	switch a.PeriodicityUnit {
	case "second", "seconds":
		scheduler.Every(a.PeriodicityValue).Seconds()
	case "minute", "minutes":
		scheduler.Every(a.PeriodicityValue).Minutes()
	case "hour", "hours":
		scheduler.Every(a.PeriodicityValue).Hours()
	case "day", "days":
		scheduler.Every(a.PeriodicityValue).Days()
	case "week", "weeks":
		scheduler.Every(a.PeriodicityValue).Weeks()
	default:
		err = errors.New("the access PeriodicityUnit cannot be translated to a scheduler")
	}

	return
}
