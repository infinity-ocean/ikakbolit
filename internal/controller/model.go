package controller

import (
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/model"
)

type Response struct {
	Schedule_id string `json:"schedule_id"`
}

type Schedule struct {
	ID          int
	UserID      string        `json:"user_id"`
	CureName    string        `json:"cure_name"`
	DosesPerDay int           `json:"doses_per_day"`
	Duration    int           `json:"duration"`
	CreatedAt   time.Time     `json:"created_at"`

	DayStart    time.Time     `json:"-"`
	DayFinish   time.Time     `json:"-"`
	
	Intakes     []string      `json:"intakes"`
}

func toModelSchedule(s Schedule) model.Schedule {
	return model.Schedule{
		ID:          s.ID,
		UserID:      s.UserID,
		CureName:    s.CureName,
		DosesPerDay: s.DosesPerDay,
		Duration:    s.Duration,
		CreatedAt:   s.CreatedAt,
	}
}

type SchedulesInWindow struct {
	Schedules []Schedule `json:"schedules"`
}
