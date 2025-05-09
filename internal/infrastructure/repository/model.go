package repository

import "time"

type Schedule struct {
	ID           int       `db:"id"`
	UserID       int       `db:"fk_user_id"`
	CureName     string    `db:"cure_name"`
	DosesPerDay  int       `db:"frequency"`
	DurationDays int       `db:"duration"`
	CreatedAt    time.Time `db:"created_at"`

	DayStart     time.Time
	DayFinish    time.Time
	
	Intakes      []string
}
