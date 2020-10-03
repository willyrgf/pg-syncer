package syncer

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron"
)

// GetScheduler translate the syncersaccess
func (a *Access) GetScheduler() (scheduler *gocron.Scheduler, err error) {
	scheduler = gocron.NewScheduler(time.Local)
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
