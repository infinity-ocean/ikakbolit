package service

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
	"github.com/stretchr/testify/require"
)

type stubRepo struct {
	schedules []entity.Schedule
	selectErr error
}

func (r *stubRepo) InsertSchedule(_ entity.Schedule) (int, error) { return 0, nil }
func (r *stubRepo) SelectSchedules(_ int) ([]entity.Schedule, error) { return r.schedules, r.selectErr }
func (r *stubRepo) SelectSchedule(_ int, _ int) (entity.Schedule, error) { return entity.Schedule{}, nil }

func TestCalculateIntakeTimes(t *testing.T) {
	log := slog.Default()
	// positive cases
	svc := New(&stubRepo{}, log, config.Config{})
	positive := []struct{ doses int; expect []string }{
		{1, []string{"08:00"}},
		{2, []string{"08:00", "22:00"}},
		{3, []string{"08:00", "15:00", "22:00"}},
		{5, []string{"08:00", "11:30", "15:00", "18:30", "22:00"}},
		{4, []string{"08:00", "12:45", "17:30", "22:00"}},
	}
	for _, tc := range positive {
		t.Run("pos_"+strconv.Itoa(tc.doses), func(t *testing.T) {
			ints, err := svc.CalculateIntakeTimes(context.Background(), tc.doses)
			require.NoError(t, err)
			require.Equal(t, tc.expect, ints)
		})
	}

	t.Run("zero_doses", func(t *testing.T) {
		ints, err := svc.CalculateIntakeTimes(context.Background(), 0)
		require.NoError(t, err)
		require.Empty(t, ints)
	})

	panicCases := []struct{ name string; doses int }{
		{"negative_doses", -1},
		{"overflow_range", -1000},
	}
	for _, tc := range panicCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Panics(t, func() {
				svc.CalculateIntakeTimes(context.Background(), tc.doses)
			})
		})
	}
}

func TestGetNextTakings(t *testing.T) {
	now := time.Now()
	cfg := config.Config{Options: config.Options{CureScheduleWindowMinutes: 60 * 24}}
	log := slog.Default()
	testCases := []struct {
		name      string
		repo      *stubRepo
		userID    int
		expectErr bool
		expectN   int
	}{
		{"repo_error", &stubRepo{selectErr: errors.New("db fail")}, 1, true, 0},
		{"no_schedules", &stubRepo{schedules: []entity.Schedule{}}, 1, false, 0},
		{"expired_all", &stubRepo{schedules: []entity.Schedule{{ID: 1, UserID: 1, DosesPerDay: 1, DurationDays: 0, CreatedAt: now.Add(-24 * time.Hour)}}}, 1, false, 0},
		{"past_intake_today", &stubRepo{schedules: []entity.Schedule{{ID: 3, UserID: 1, DosesPerDay: 1, DurationDays: 1, CreatedAt: now.Add(-1 * time.Minute)}}}, 1, false, 1},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.repo
			svc := &Service{repo: repo, log: log, cfg: cfg}
			res, err := svc.GetNextTakings(context.Background(), tc.userID)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, res, tc.expectN)
			}
		})
	}
}
