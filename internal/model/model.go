package model

import (
	"os"
	"time"
)

type ScheduleRequest struct {
	ID        int           `json:"id"`
	UserID    string        `json:"user_id"`
	CureName  string        `json:"cure_name"`
	Frequency time.Duration `json:"frequency"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"created_at"`
}

type ScheduleDB struct {
	ID        int           `db:"id"`
	UserID    string        `db:"fk_user_id"`
	CureName  string        `db:"cure_name"`
	Frequency time.Duration `db:"frequency"`
	Duration  time.Duration `db:"duration"`
	CreatedAt time.Time     `db:"created_at"`

	DayStart  time.Time
	DayFinish time.Time
}

func ToScheduleDB(req ScheduleRequest) (ScheduleDB, error) {
	dayStartStr := os.Getenv("DAY_START")
	dayStart, err := time.Parse("15:04", dayStartStr)
	if err != nil {
		return ScheduleDB{}, err
	}

	dayFinishStr := os.Getenv("DAY_FINISH")
	dayFinish, err := time.Parse("15:04", dayFinishStr)
	if err != nil {
		return ScheduleDB{}, err
	}

	return ScheduleDB{
		ID:        req.ID,
		UserID:    req.UserID,
		CureName:  req.CureName,
		Frequency: req.Frequency,
		Duration:  req.Duration,
		CreatedAt: req.CreatedAt,
		DayStart:  dayStart,
		DayFinish: dayFinish,
	}, nil
}
