package model

import (
	"time"
)

type Schedule struct {
	ID          int
	UserID      int
	CureName    string
	DosesPerDay int
	Duration    int
	CreatedAt   time.Time
	
	DayStart    time.Time
	DayFinish   time.Time

	Intakes     []string
}

