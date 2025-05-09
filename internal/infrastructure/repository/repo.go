package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) InsertSchedule(sched entity.Schedule) (int, error) {
	ctx := context.Background()
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
			return 0, fmt.Errorf("failed to accquire connection pool: %w", err)
	}
	defer conn.Release()

	var id int
	query := `INSERT INTO scheduled (user_id, cure_name, doses_per_day, duration_days) 
			  VALUES ($1, $2, $3, $4) RETURNING id;`

	err = conn.QueryRow(ctx, query, sched.UserID, sched.CureName, sched.DosesPerDay, sched.DurationDays).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert schedule: %w", err)
	}

	return id, nil
}

func (r *Repo) SelectSchedules(userID int) ([]entity.Schedule, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to accquire connection pool: %w", err)
	}
	defer conn.Release()

	sql := `
		SELECT id, user_id, cure_name, doses_per_day, duration_days, created_at 
		FROM scheduled 
		WHERE user_id = $1
	`
	rows, err := conn.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select schedules: %w", err)
	}
	defer rows.Close()

	var schedules []entity.Schedule
	for rows.Next() {
		var s entity.Schedule
		var duration int

		if err := rows.Scan(&s.ID, &s.UserID, &s.CureName, &s.DosesPerDay, &duration, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		s.DurationDays = duration
		schedules = append(schedules, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return schedules, nil
}

func (r *Repo) SelectSchedule(userID int, schedID int) (entity.Schedule, error) {
	conn, err := r.pool.Acquire(context.Background())
	if err != nil {
		return entity.Schedule{}, fmt.Errorf("failed to accquire connection pool: %w", err)
	}
	defer conn.Release()

	sql := `SELECT id, user_id, cure_name, doses_per_day, duration_days, created_at 
			FROM scheduled 
			WHERE user_id = $1 AND id = $2`

	row := conn.QueryRow(context.Background(), sql, userID, schedID)

	var schedule entity.Schedule
	err = row.Scan(
		&schedule.ID,
		&schedule.UserID,
		&schedule.CureName,
		&schedule.DosesPerDay,
		&schedule.DurationDays,
		&schedule.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Schedule{}, fmt.Errorf("schedule not found: %w", err)
		}
		return entity.Schedule{}, fmt.Errorf("failed to scan schedule: %w", err)
	}

	return schedule, nil
}
