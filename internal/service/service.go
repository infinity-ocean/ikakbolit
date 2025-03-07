package service

import "github.com/infinity-ocean/ikakbolit/internal/model"

type service struct {
	repo repo
}

type repo interface {
	InsertSchedule(model.ScheduleDB) (int, error)
}

func New(repo repo) *service {
	return &service{repo: repo}
}

func (s *service) AddSchedule(scheduleReq model.ScheduleRequest) (int, error) {
	scheduleDB, err := ToScheduleDB(scheduleReq)
	if err != nil {
		return 0, err
	}

	return s.repo.InsertSchedule(scheduleDB)
}