package model

import (
	"time"
)

type ScheduleRequest struct {
	ID        int
	UserID    string
	CureName  string
	Frequency time.Duration
	Duration  time.Duration
	CreatedAt time.Time
}

type ScheduleDB struct {
	ID        int          
	UserID    string       
	CureName  string       
	Frequency time.Duration
	Duration  time.Duration
	CreatedAt time.Time    

	DayStart  time.Time
	DayFinish time.Time
}


