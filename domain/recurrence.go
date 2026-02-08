package domain

import (
	"time"
)

type Frequency string

const (
	FreqDaily   Frequency = "Daily"
	FreqWeekly  Frequency = "Weekly"
	FreqMonthly Frequency = "Monthly"
)

type Recurrence struct {
	id         int64
	frequency  Frequency
	interval   int
	dayOfWeek  time.Weekday
	dayOfMonth int
	startDate  time.Time
	endDate    time.Time
}

func NewRecurrence(
	id int64,
	frequencyString string,
	interval int,
	dayOfWeek time.Weekday,
	dayOfMonth int,
	startDate time.Time,
	endDate time.Time,
) (*Recurrence, *FreqError) {
	var frequency Frequency
	switch frequencyString {
	case "Daily":
		frequency = FreqDaily
	case "Weekly":
		frequency = FreqWeekly
	case "Monthly":
		frequency = FreqMonthly
	default:
		return nil, &FreqError{frequencyString, "Failed to create new Recurrence"}
	}
	return &Recurrence{
		id,
		frequency,
		interval,
		dayOfWeek,
		dayOfMonth,
		startDate,
		endDate,
	}, nil
}
