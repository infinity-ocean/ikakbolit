package entity

import (
	"time"
)

type Schedule struct {
	ID           int
	UserID       int
	CureName     string
	DosesPerDay  int
	DurationDays int
	CreatedAt    time.Time

	DayStart     time.Time
	DayFinish    time.Time
	
	Intakes      []string
}
