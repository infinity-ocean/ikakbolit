package repo

import "time"

type Schedule struct {
	ID          int           `db:"id"`
	UserID      string        `db:"fk_user_id"`
	CureName    string        `db:"cure_name"`
	DosesPerDay int           `db:"frequency"`
	Duration    int 		  `db:"duration"`
	CreatedAt   time.Time     `db:"created_at"`

	DayStart    time.Time
	DayFinish   time.Time

	Intakes     []string
}
