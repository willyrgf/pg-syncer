package syncer

import (
	"errors"

	"github.com/jasonlvhit/gocron"
)

// GetScheduler translate the syncersaccess
func (a *Access) GetScheduler() (scheduler *gocron.Job, err error) {
	switch a.PeriodicityUnit {
	case "second", "seconds":
		scheduler = gocron.Every(a.PeriodicityValue).Seconds()
	case "minute", "minutes":
		scheduler = gocron.Every(a.PeriodicityValue).Minutes()
	case "hour", "hours":
		scheduler = gocron.Every(a.PeriodicityValue).Hours()
	case "day", "days":
		scheduler = gocron.Every(a.PeriodicityValue).Days()
	case "week", "weeks":
		scheduler = gocron.Every(a.PeriodicityValue).Weeks()
	default:
		err = errors.New("the access PeriodicityUnit cannot be translated to a scheduler")
	}

	if err == nil {
		err = scheduler.Err()
	}

	return
}
