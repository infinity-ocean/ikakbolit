package model

import (
	"time"
)

type Schedule struct {
	ID        int          
	UserID    string       
	CureName  string       
	DosesPerDay int
	Duration  time.Duration
	CreatedAt time.Time    

	DayStart  time.Time
	DayFinish time.Time

	Intakes     []string
}