package repo

import "time"

type ScheduleDB struct {
	ID                     int           `db:"id"`
	UserID                 string        `db:"fk_user_id"`
	CureName               string        `db:"cure_name"`
	DosesPerDayDosesPerDay int           `db:"frequency"`
	Duration               time.Duration `db:"duration"`
	CreatedAt              time.Time     `db:"created_at"`

	DayStart  time.Time
	DayFinish time.Time
}
