package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/domain/entity"
	"github.com/infinity-ocean/ikakbolit/internal/infrastructure/repository"
	"github.com/infinity-ocean/ikakbolit/pkg/errcodes"
)

type Service struct {
	repo repo
	log  *slog.Logger
	cfg  config.Config
}

type repo interface {
	InsertSchedule(entity.Schedule) (int, error)
	SelectSchedules(int) ([]entity.Schedule, error)
	SelectSchedule(int, int) (entity.Schedule, error)
}

func New(repo repo, log *slog.Logger, cfg config.Config) *Service {
	return &Service{repo: repo, log: log, cfg: cfg}
}

func (s *Service) AddSchedule(ctx context.Context, schedule entity.Schedule) (int, error) {
	if schedule.DosesPerDay < 1 || schedule.DosesPerDay > 12 {
		return 0, fmt.Errorf("doses per day must be between 1 and 12: %w", errcodes.ErrBadRequest)
	}

	log := s.logger(ctx)
	log.Debug("starting to add schedule", "user_id", schedule.UserID)

	id, err := s.repo.InsertSchedule(schedule)
	if err != nil {
		log.Error("failed to insert schedule", "error", err)
		return 0, repository.MapSQLError(err)
	}

	log.Info("schedule added successfully", "schedule_id", id)
	return id, nil
}

func (s *Service) GetScheduleIDs(ctx context.Context, userID int) ([]int, error) {
	log := s.logger(ctx)
	log.Debug("starting to get schedule IDs", "user_id", userID)

	schedules, err := s.repo.SelectSchedules(userID)
	if err != nil {
		log.Error("failed to select schedules", "user_id", userID, "error", err)
		return nil, repository.MapSQLError(err)
	}

	idSlice := make([]int, 0, len(schedules))
	for _, schedule := range schedules {
		idSlice = append(idSlice, schedule.ID)
	}

	log.Info("retrieved schedule IDs", "user_id", userID, "count", len(idSlice))
	return idSlice, nil
}

func (s *Service) GetScheduleWithIntake(ctx context.Context, userID int, scheduleID int) (entity.Schedule, error) {
	log := s.logger(ctx)
	log.Debug("starting to get schedule with intake", "user_id", userID, "schedule_id", scheduleID)

	schedule, err := s.repo.SelectSchedule(userID, scheduleID)
	if err != nil {
		log.Error("failed to select schedule", "user_id", userID, "schedule_id", scheduleID, "error", err)
		return entity.Schedule{}, repository.MapSQLError(err)
	}

	if schedule.ID == 0 {
		log.Error("invalid schedule returned (ID is 0)", "user_id", userID, "schedule_id", scheduleID)
		return entity.Schedule{}, fmt.Errorf("schedule_id is 0: %w", err)
	}

	intakes, err := s.CalculateIntakeTimes(ctx, schedule.DosesPerDay)
	if err != nil {
		log.Error("failed to calculate intake times", "user_id", userID, "schedule_id", scheduleID, "error", err)
		return entity.Schedule{}, fmt.Errorf("error in schedule calculating: %w", err)
	}
	schedule.Intakes = intakes

	log.Info("schedule with intake retrieved", "user_id", userID, "schedule_id", scheduleID, "intake_count", len(intakes))
	return schedule, nil
}

func (s *Service) GetNextTakings(ctx context.Context, userID int) ([]entity.Schedule, error) {
	log := s.logger(ctx)
	log.Debug("starting to compute next takings", "user_id", userID)

	schedules, err := s.repo.SelectSchedules(userID)
	if err != nil {
		log.Error("failed to select schedules", "user_id", userID, "error", err)
		return nil, repository.MapSQLError(err)
	}
	if len(schedules) == 0 {
		log.Info("no schedules found for user", "user_id", userID)
		return nil, nil
	}

	windowMin := s.cfg.Options.CureScheduleWindowMinutes
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
			log.Error("failed to calculate intake times", "user_id", userID, "schedule_id", sch.ID, "error", err)
			return nil, fmt.Errorf("incorrect .env DAY_START or DAY_FINISH: %w", err)
		}
		sch.Intakes = times

		for _, timeStr := range times {
			intakeTime, err := time.ParseInLocation("15:04", timeStr, now.Location())
			if err != nil {
				log.Error("failed to parse intake time", "user_id", userID, "schedule_id", sch.ID, "time_str", timeStr, "error", err)
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

	log.Info("computed next takings", "user_id", userID, "next_count", len(result))
	return result, nil
}

func (s *Service) CalculateIntakeTimes(ctx context.Context, dosesPerDay int) ([]string, error) {
	log := s.logger(ctx)
	log.Debug("Starting calculation of intake times", "doses_per_day", dosesPerDay)

	dayStartStr := "08:00"
	dayFinishStr := "22:00"
	log.Debug("Using fixed day period for intake calculation", "day_start", dayStartStr, "day_finish", dayFinishStr)

	dayStart, err := time.Parse("15:04", dayStartStr)
	if err != nil {
		log.Error("Failed to parse day start", "day_start", dayStartStr, "error", err)
		return nil, fmt.Errorf("failed to parse day start: %w", err)
	}

	dayFinish, err := time.Parse("15:04", dayFinishStr)
	if err != nil {
		log.Error("Failed to parse day finish", "day_finish", dayFinishStr, "error", err)
		return nil, fmt.Errorf("failed to parse day finish: %w", err)
	}

	log.Debug("Parsed day period boundaries", "day_start_time", dayStart.Format("15:04"), "day_finish_time", dayFinish.Format("15:04"))

	if dosesPerDay == 1 {
		intake := dayStartStr
		log.Info("Single daily intake time", "intake_time", intake)
		return []string{intake}, nil
	}

	totalDuration := dayFinish.Sub(dayStart)
	interval := totalDuration / time.Duration(dosesPerDay-1)
	log.Debug("Computed interval between doses", "interval_minutes", interval.Minutes())

	rawIntakes := make([]time.Time, dosesPerDay)
	for i := range dosesPerDay {
		rawIntakes[i] = dayStart.Add(time.Duration(i) * interval)
	}
	log.Debug("Raw intake times calculated", "raw_intakes", rawIntakes)

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
		log.Debug("Rounded intake time", "index", i, "from", rawIntakes[i].Format("15:04"), "to", intakes[i])
	}

	log.Info("Final intake times after rounding", "intakes", intakes)
	return intakes, nil
}

func (s *Service) logger(ctx context.Context) *slog.Logger {
	reqID := GetReqID(ctx)
	if reqID == "" {
		return s.log
	}
	return s.log.With("request_id", reqID)
}

type ctxKey string

const requestIDKey ctxKey = "request_id"

func GetReqID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
