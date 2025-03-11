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
	ctx := context.Background()
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
			return 0, fmt.Errorf("failed to accquire connection pool: %w", model.ErrInternalServerError)
	}
	defer conn.Release()

	var id int
	query := `INSERT INTO scheduled (user_id, cure_name, doses_per_day, duration) 
			  VALUES ($1, $2, $3, $4) RETURNING id;`

	err = conn.QueryRow(ctx, query, sched.UserID, sched.CureName, sched.DosesPerDay, sched.Duration).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert schedule: %w", model.ErrInternalServerError)
	}

	return id, nil
}

func (r *repo) SelectSchedules(userID int) ([]model.Schedule, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to accquire connection pool: %w", model.ErrInternalServerError)
	}
	defer conn.Release()

	sql := `
		SELECT id, user_id, cure_name, doses_per_day, duration, created_at 
		FROM scheduled 
		WHERE user_id = $1
	`
	rows, err := conn.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select schedules: %w", model.ErrInternalServerError)
	}
	defer rows.Close()

	var schedules []model.Schedule
	for rows.Next() {
		var s model.Schedule
		var duration int

		if err := rows.Scan(&s.ID, &s.UserID, &s.CureName, &s.DosesPerDay, &duration, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", model.ErrInternalServerError)
		}

		s.Duration = duration
		schedules = append(schedules, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", model.ErrInternalServerError)
	}

	return schedules, nil
}

func (r *repo) SelectSchedule(userID int, schedID int) (model.Schedule, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return model.Schedule{}, fmt.Errorf("failed to accquire connection pool: %w", model.ErrInternalServerError)
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
			return model.Schedule{}, nil
		}
		return model.Schedule{}, fmt.Errorf("failed to scan schedule: %w", model.ErrInternalServerError)
	}
	return schedule, nil
}