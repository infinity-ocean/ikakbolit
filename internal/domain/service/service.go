package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/infinity-ocean/ikakbolit/internal/dto"
	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type service struct {
	repo   repo
	log *slog.Logger
}

type repo interface {
	InsertSchedule(entity.Schedule) (int, error)
	SelectSchedules(int) ([]entity.Schedule, error)
	SelectSchedule(int, int) (entity.Schedule, error)
}

func New(repo repo, log *slog.Logger) *service {
	return &service{repo: repo, log: log}
}

func (s *service) AddSchedule(ctx context.Context, schedule entity.Schedule) (int, error) {
	if schedule.DosesPerDay < 1 || schedule.DosesPerDay > 12 {
		return 0, fmt.Errorf("doses per day must be between 1 and 12: %w", dto.ErrBadRequest)
	}

	reqID := middleware.GetReqID(ctx)
	s.log.Info("starting to add schedule", "request_id", reqID, "user_id", schedule.UserID)

	id, err := s.repo.InsertSchedule(schedule)
	if err != nil {
		s.log.Error("failed to insert schedule", "request_id", reqID, "error", err)
		return 0, MapSQLErrorToDTO(err)
	}

	s.log.Info("schedule added successfully", "request_id", reqID, "schedule_id", id)
	return id, nil
}

func (s *service) GetScheduleIDs(ctx context.Context, userID int) ([]int, error) {
    reqID := middleware.GetReqID(ctx)
    s.log.Info("starting to get schedule IDs", "request_id", reqID, "user_id", userID)

    schedules, err := s.repo.SelectSchedules(userID)
    if err != nil {
        s.log.Error("failed to select schedules", "request_id", reqID, "user_id", userID, "error", err)
        return nil, MapSQLErrorToDTO(err)
    }

    idSlice := make([]int, 0, len(schedules))
    for _, schedule := range schedules {
        idSlice = append(idSlice, schedule.ID)
    }

    s.log.Info("retrieved schedule IDs", "request_id", reqID, "user_id", userID, "count", len(idSlice))
    return idSlice, nil
}

func (s *service) GetScheduleWithIntake(ctx context.Context, userID int, scheduleID int) (entity.Schedule, error) {
    reqID := middleware.GetReqID(ctx)
    s.log.Info("starting to get schedule with intake", "request_id", reqID, "user_id", userID, "schedule_id", scheduleID)

    schedule, err := s.repo.SelectSchedule(userID, scheduleID)
    if err != nil {
        s.log.Error("failed to select schedule", "request_id", reqID, "user_id", userID, "schedule_id", scheduleID, "error", err)
        return entity.Schedule{}, MapSQLErrorToDTO(err)
    }

    if schedule.ID == 0 {
        s.log.Error("invalid schedule returned (ID is 0)", "request_id", reqID, "user_id", userID, "schedule_id", scheduleID)
        return entity.Schedule{}, fmt.Errorf("schedule_id is 0: %w", err)
    }

    intakes, err := s.CalculateIntakeTimes(ctx, schedule.DosesPerDay)
    if err != nil {
        s.log.Error("failed to calculate intake times", "request_id", reqID, "user_id", userID, "schedule_id", scheduleID, "error", err)
        return entity.Schedule{}, fmt.Errorf("error in schedule calculating: %w", err)
    }
    schedule.Intakes = intakes

    s.log.Info("schedule with intake retrieved", "request_id", reqID, "user_id", userID, "schedule_id", scheduleID, "intake_count", len(intakes))
    return schedule, nil
}

func (s *service) GetNextTakings(ctx context.Context, userID int) ([]entity.Schedule, error) {
    reqID := middleware.GetReqID(ctx)
    s.log.Info("starting to compute next takings", "request_id", reqID, "user_id", userID)

    schedules, err := s.repo.SelectSchedules(userID)
    if err != nil {
        s.log.Error("failed to select schedules", "request_id", reqID, "user_id", userID, "error", err)
        return nil, MapSQLErrorToDTO(err)
    }
    if len(schedules) == 0 {
        s.log.Info("no schedules found for user", "request_id", reqID, "user_id", userID)
        return nil, nil
    }

    windowMinStr := os.Getenv("CURE_SCHEDULE_WINDOW_MIN")
    windowMin, err := strconv.Atoi(windowMinStr)
    if err != nil {
        s.log.Error("invalid CURE_SCHEDULE_WINDOW_MIN env var", "request_id", reqID, "value", windowMinStr, "error", err)
        return nil, fmt.Errorf("env CURE_SCHEDULE_WINDOW_MIN is not int: %w", err)
    }
    windowDuration := time.Duration(windowMin) * time.Minute

    now := time.Now()
    var result []entity.Schedule

    for _, sch := range schedules {
        scheduleEnd := sch.CreatedAt.Add(time.Duration(sch.DurationDays) * 24 * time.Hour)
        if now.After(scheduleEnd) {
            continue
        }

        times, err := s.CalculateIntakeTimes(ctx, sch.DosesPerDay)
        if err != nil {
            s.log.Error("failed to calculate intake times", "request_id", reqID, "user_id", userID, "schedule_id", sch.ID, "error", err)
            return nil, fmt.Errorf("incorrect .env DAY_START or DAY_FINISH: %w", err)
        }
        sch.Intakes = times

        for _, timeStr := range times {
            intakeTime, err := time.ParseInLocation("15:04", timeStr, now.Location())
            if err != nil {
                s.log.Error("failed to parse intake time", "request_id", reqID, "user_id", userID, "schedule_id", sch.ID, "time_str", timeStr, "error", err)
                return nil, fmt.Errorf("failed to parse intake time %s: %w", timeStr, err)
            }
            intakeTime = time.Date(now.Year(), now.Month(), now.Day(), intakeTime.Hour(), intakeTime.Minute(), 0, 0, now.Location())
            if intakeTime.Before(now) {
                intakeTime = intakeTime.Add(24 * time.Hour)
            }
            if intakeTime.After(now) && intakeTime.Before(now.Add(windowDuration)) {
                result = append(result, sch)
                break
            }
        }
    }

    s.log.Info("computed next takings", "request_id", reqID, "user_id", userID, "next_count", len(result))
    return result, nil
}

func (s *service) CalculateIntakeTimes(ctx context.Context, dosesPerDay int) ([]string, error) {
      reqID := middleware.GetReqID(ctx)
      s.log.Info("Starting calculation of intake times", "request_id", reqID, "doses_per_day", dosesPerDay)

      dayStartStr := "08:00"
      dayFinishStr := "22:00"
      s.log.Info("Using fixed day period for intake calculation", "request_id", reqID, "day_start", dayStartStr, "day_finish", dayFinishStr)

      dayStart, err := time.Parse("15:04", dayStartStr)
      if err != nil {
          s.log.Error("Failed to parse day start", "request_id", reqID, "day_start", dayStartStr, "error", err)
          return nil, fmt.Errorf("failed to parse day start: %w", err)
      }

      dayFinish, err := time.Parse("15:04", dayFinishStr)
      if err != nil {
          s.log.Error("Failed to parse day finish", "request_id", reqID, "day_finish", dayFinishStr, "error", err)
          return nil, fmt.Errorf("failed to parse day finish: %w", err)
      }

      s.log.Debug("Parsed day period boundaries", "request_id", reqID, "day_start_time", dayStart.Format("15:04"), "day_finish_time", dayFinish.Format("15:04"))

      if dosesPerDay == 1 {
          intake := dayStartStr
          s.log.Info("Single daily intake time", "request_id", reqID, "intake_time", intake)
          return []string{intake}, nil
      }

      totalDuration := dayFinish.Sub(dayStart)
      interval := totalDuration / time.Duration(dosesPerDay-1)
      s.log.Debug("Computed interval between doses", "request_id", reqID, "interval_minutes", interval.Minutes())

      rawIntakes := make([]time.Time, dosesPerDay)
      for i := range dosesPerDay {
          rawIntakes[i] = dayStart.Add(time.Duration(i) * interval)
      }
      s.log.Debug("Raw intake times calculated", "request_id", reqID, "raw_intakes", rawIntakes)

      intakes := make([]string, dosesPerDay)
      for i, t := range rawIntakes {
          minute := t.Minute()
          if minute%15 != 0 {
              minute += 15 - (minute % 15)
              if minute == 60 {
                  t = t.Add(time.Hour)
                  minute = 0
              }
          }
          rounded := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 0, 0, t.Location())
          intakes[i] = rounded.Format("15:04")
          s.log.Debug("Rounded intake time", "request_id", reqID, "index", i, "from", rawIntakes[i].Format("15:04"), "to", intakes[i])
      }

      s.log.Info("Final intake times after rounding", "request_id", reqID, "intakes", intakes)
      return intakes, nil
}

func MapSQLErrorToDTO(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %v", dto.ErrNotFound, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23502", "23503", "23505", "23514", "22001", "22002", "22003", "22P02", "22007", "42804":
			return fmt.Errorf("%w: %s", dto.ErrBadRequest, pgErr.Message)
		}
	}

	return err
}
