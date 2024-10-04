package service

import (
	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/robfig/cron/v3"
)

type scheduler struct{}

func NewScheduler() *scheduler {
	cron.New()
	return &scheduler{}
}

func (s *scheduler) Register(cron string, name string, plain *domain.Plan) error {
	return nil
}

func (s *scheduler) Start() error {
	return nil
}
