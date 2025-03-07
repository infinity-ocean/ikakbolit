package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/infinity-ocean/ikakbolit/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repo {
	return &repo{pool: pool}
}

func (r *repo) InsertSchedule(sched model.Schedule) (int, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	var id int
	sql := `INSERT INTO scheduled (user_id, cure_name, doses_per_day, duration, created_at) 
			VALUES ($1, $2, $3, $4, $5)						
			RETURNING id;`
	err = conn.QueryRow(
		context.Background(),
		sql, sched.UserID,
		sched.CureName,
		sched.DosesPerDay,
		sched.Duration,
		sched.CreatedAt).
		Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) SelectSchedules(userID int) ([]int, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `SELECT id FROM scheduled WHERE user_id = $1`
	rows, err := conn.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduleIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		scheduleIDs = append(scheduleIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return scheduleIDs, nil
}

func (r *repo) SelectSchedule(userID int, schedID int) (model.Schedule, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return model.Schedule{}, err
	}
	defer conn.Release()

	sql := `SELECT id, user_id, cure_name, doses_per_day, duration, created_at 
			FROM scheduled 
			WHERE user_id = $1 AND id = $2`

	row := conn.QueryRow(context.Background(), sql, userID, schedID)

	var schedule model.Schedule
	err = row.Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.CureName,
		&schedule.DosesPerDay,
		&schedule.Duration,
		&schedule.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Schedule{}, fmt.Errorf("schedule not found for user_id=%d, schedID=%d", userID, schedID)
		}
		return model.Schedule{}, err
	}
	return schedule, nil
}