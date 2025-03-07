package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/model"
)

type service struct {
	repo repo
}

type repo interface {
	InsertSchedule(model.Schedule) (int, error)
	SelectSchedules(int) ([]int, error)
	SelectSchedule(int, int) (model.Schedule, error)
}

func New(repo repo) *service {
	return &service{repo: repo}
}

func (s *service) AddSchedule(schedule model.Schedule) (int, error) {
	if schedule.DosesPerDay < 1 || schedule.DosesPerDay > 24 {
		return 0, errors.New("doses per day is <1 or >24")
	}

	return s.repo.InsertSchedule(schedule)
}

func (s *service) GetSchedules(userID int) ([]int, error) {
	return s.repo.SelectSchedules(userID)
}

func (s *service) GetSchedule(userID int, scheduleID int) (model.Schedule, []string, error) {
	schedule, err := s.repo.SelectSchedule(userID, scheduleID)
	if err != nil {
		return model.Schedule{}, nil, err
	}
	intakes, err := CalculateIntakeTimes(schedule.DosesPerDay)
	if err != nil {
		return model.Schedule{}, nil, err
	}

	return schedule, intakes, nil
}


func CalculateIntakeTimes(dosesPerDay int) ([]string, error) {
	dayStartStr := os.Getenv("DAY_START")
	dayFinishStr := os.Getenv("DAY_FINISH`")

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