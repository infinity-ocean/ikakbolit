package controller

import (
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/model"
)

type Response struct {
	Schedule_id string `json:"schedule_id"`
}

type ScheduleRequest struct {
	ID          int
	UserID      string        `json:"user_id"`
	CureName    string        `json:"cure_name"`
	DosesPerDay int           `json:"DosesPerDay"`
	Duration    time.Duration `json:"duration"`
	CreatedAt   time.Time     `json:"created_at"`
}

func toModelSchedule(s ScheduleRequest) model.Schedule {
	return model.Schedule{
		ID:          s.ID,
		UserID:      s.UserID,
		CureName:    s.CureName,
		DosesPerDay: s.DosesPerDay,
		Duration:    s.Duration,
		CreatedAt:   s.CreatedAt,
	}
}

type ScheduleWithIntakes struct {
	ID          int           `id:"user_id"`
	UserID      string        `json:"user_id"`
	CureName    string        `json:"cure_name"`
	DosesPerDay int           `json:"DosesPerDay"`
	Duration    time.Duration `json:"duration"`
	CreatedAt   time.Time     `json:"created_at"`
	Intakes     []string      `json:"intakes"`
}

func toScheduleWithIntakes(s model.Schedule, intakes []string) ScheduleWithIntakes {
	return ScheduleWithIntakes{
		ID:          s.ID,
		UserID:      s.UserID,
		CureName:    s.CureName,
		DosesPerDay: s.DosesPerDay,
		Duration:    s.Duration,
		CreatedAt:   s.CreatedAt,
		Intakes:     intakes,
	}
}
