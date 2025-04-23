package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/model"
)

type service struct {
	repo repo
}

// GetSchedule implements controller.IkakbolitService.
func (s *service) GetSchedule(int, int) (model.Schedule, error) {
	panic("unimplemented")
}

type repo interface {
	InsertSchedule(model.Schedule) (int, error)
	SelectSchedules(int) ([]model.Schedule, error)
	SelectSchedule(int, int) (model.Schedule, error)
}

func New(repo repo) *service {
	return &service{repo: repo}
}

func (s *service) AddSchedule(schedule model.Schedule) (int, error) {
	if schedule.DosesPerDay < 1 || schedule.DosesPerDay > 24 {
		return 0, fmt.Errorf("doses per day must be between 1 and 24: %w", model.ErrBadRequest)
	}
	return s.repo.InsertSchedule(schedule)
}

func (s *service) GetScheduleIDs(userID int) ([]int, error) {
	schedules, err := s.repo.SelectSchedules(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule IDs: %w", err)
	}
	idSlice := make([]int, 0, len(schedules))
	for _, schedule := range schedules {
		idSlice = append(idSlice, schedule.ID)
	}
	return idSlice, nil
}

func (s *service) GetScheduleWithIntake(userID int, scheduleID int) (model.Schedule, error) {
	schedule, err := s.repo.SelectSchedule(userID, scheduleID)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("failed to get schedule: %w", err)
	}

	if schedule.ID == 0 {
		return model.Schedule{}, fmt.Errorf("no schedule for user %d with scheduleID %d: %w", userID, scheduleID, model.ErrNoContent)
	}

	intakes, err := CalculateIntakeTimes(schedule.DosesPerDay)
	if err != nil {
		return model.Schedule{}, fmt.Errorf("incorrect .env DAY_START or DAY_FINISH: %w", model.ErrInternalServerError)
	}
	schedule.Intakes = intakes
	return schedule, nil
}

func (s *service) GetNextTakings(userID int) ([]model.Schedule, error) {
	schedules, err := s.repo.SelectSchedules(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}

	if len(schedules) == 0 {
		return nil, fmt.Errorf("no schedules for user %d: %w", userID, model.ErrNoContent)
	}

	window, err := strconv.Atoi(os.Getenv("CURE_SCHEDULE_WINDOW_MIN"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse CURE_SCHEDULE_WINDOW_MIN: %w", model.ErrInternalServerError)
	}
	windowDuration := time.Duration(window) * time.Minute

	now := time.Now()
	var result []model.Schedule

	for _, schedule := range schedules {
		scheduleEnd := schedule.CreatedAt.Add(time.Duration(schedule.DurationDays) * 24 * time.Hour)
		if now.After(scheduleEnd) {
			continue
		}

		times, err := CalculateIntakeTimes(schedule.DosesPerDay)
		if err != nil {
			return nil, fmt.Errorf("incorrect .env DAY_START or DAY_FINISH: %w", model.ErrInternalServerError)
		}
		schedule.Intakes = times

		for _, timeStr := range times {
			intakeTime, err := time.ParseInLocation("15:04", timeStr, now.Location())
			if err != nil {
				return nil, fmt.Errorf("failed to parse intake time %s: %w", timeStr, model.ErrInternalServerError)
			}

			intakeTime = time.Date(now.Year(), now.Month(), now.Day(), intakeTime.Hour(), intakeTime.Minute(), 0, 0, now.Location())

			if intakeTime.Before(now) {
				intakeTime = intakeTime.Add(24 * time.Hour)
			}

			if intakeTime.After(now) && intakeTime.Before(now.Add(windowDuration)) {
				result = append(result, schedule)
				break
			}
		}
	}

	return result, nil
}

func CalculateIntakeTimes(dosesPerDay int) ([]string, error) {
	dayStartStr := os.Getenv("DAY_START")
	dayFinishStr := os.Getenv("DAY_FINISH")

	if dayStartStr == "" || dayFinishStr == "" {
		return nil, fmt.Errorf("environment variables DAY_START or DAY_FINISH are not set")
	}

	dayStart, err := time.Parse("15:04", dayStartStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DAY_START: %w", err)
	}

	dayFinish, err := time.Parse("15:04", dayFinishStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DAY_FINISH: %w", err)
	}

	if dosesPerDay == 1 {
		return []string{dayStart.Format("15:04")}, nil
	}

	totalDuration := dayFinish.Sub(dayStart)
	interval := totalDuration / time.Duration(dosesPerDay-1)

	intakes := make([]string, dosesPerDay)
	for i := 0; i < dosesPerDay; i++ {
		intakeTime := dayStart.Add(time.Duration(i) * interval)
		intakes[i] = intakeTime.Format("15:04")
	}

	return intakes, nil
}
