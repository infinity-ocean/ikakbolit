package rest

import (
	"time"
	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
)

// @name responseScheduleID
type ResponseScheduleID struct {
	Schedule_id string `json:"schedule_id"`
}

// swagger:model Schedule
type Schedule struct {
	ID           int
	UserID       int       `json:"user_id"`
	CureName     string    `json:"cure_name"`
	DosesPerDay  int       `json:"doses_per_day"`
	DurationDays int       `json:"duration_days"`
	CreatedAt    time.Time `json:"created_at"`

	DayStart     time.Time `json:"-"`
	DayFinish    time.Time `json:"-"`
	
	Intakes      []string  `json:"intakes"`
}

func ToModelSchedule(s Schedule) entity.Schedule {
	return entity.Schedule{
		ID:           s.ID,
		UserID:       s.UserID,
		CureName:     s.CureName,
		DosesPerDay:  s.DosesPerDay,
		DurationDays: s.DurationDays,
		CreatedAt:    s.CreatedAt,
	}
}

// swagger:model ScheduleInWindow
type SchedulesInWindow struct {
	Schedules []Schedule `json:"schedules"`
}

// swagger:model APIError
type APIError struct {
	Message string `json:"message"`
}