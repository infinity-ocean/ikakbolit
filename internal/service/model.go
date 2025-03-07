package service

import (
	"os"
	"time"

	"github.com/infinity-ocean/ikakbolit/internal/model"
)

func ToScheduleDB(req model.ScheduleRequest) (model.ScheduleDB, error) {
	dayStartStr := os.Getenv("DAY_START")
	dayStart, err := time.Parse("15:04", dayStartStr)
	if err != nil {
		return model.ScheduleDB{}, err
	}

	dayFinishStr := os.Getenv("DAY_FINISH")
	dayFinish, err := time.Parse("15:04", dayFinishStr)
	if err != nil {
		return model.ScheduleDB{}, err
	}

	return model.ScheduleDB{
		ID:        req.ID,
		UserID:    req.UserID,
		CureName:  req.CureName,
		Frequency: req.Frequency,
		Duration:  req.Duration,
		CreatedAt: req.CreatedAt,
		DayStart:  dayStart,
		DayFinish: dayFinish,
	}, nil
}