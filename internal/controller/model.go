package controller

import (
	"time"
)

type Response struct {
	Schedule_id string `json:"schedule_id"`
}

type ScheduleRequest struct {
	ID        int           `json:"id"`
	UserID    string        `json:"user_id"`
	CureName  string        `json:"cure_name"`
	Frequency time.Duration `json:"frequency"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"created_at"`
}
