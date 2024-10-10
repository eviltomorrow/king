package service

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type scheduler struct {
	cron *cron.Cron

	plans []*domain.Plan
}

func NewScheduler() *scheduler {
	return &scheduler{
		cron: cron.New(),

		plans: make([]*domain.Plan, 0, 16),
	}
}

func (s *scheduler) Register(cron string, plan *domain.Plan) {
	_, err := s.cron.AddFunc(cron, func() {
		if plan.IsCompleted() {
			return
		}

		var (
			status = domain.Ready
			err    error
		)

		if plan.Precondition != nil {
			status, err = plan.Precondition()
			if err != nil {
				zlog.Error("Check precondition failure", zap.Error(err), zap.String("name", plan.GetName()))
				return
			}

			switch status {
			case domain.Pending:
				return
			case domain.Ready:
			case domain.Completed:
				plan.SetStatus(domain.Completed)
				return
			default:
				return
			}
		}

		msg := ""
		msg, err = plan.Todo()
		if err != nil {
			zlog.Error("Todo execute failure", zap.Error(err), zap.String("name", plan.GetName()))
		} else {
			zlog.Info("Todo execute success", zap.String("name", plan.GetName()))
		}

		if plan.NotifyWithError != nil && err != nil {
			err = plan.NotifyWithError(err)
			if err != nil {
				zlog.Error("NotifyWithError failure", zap.Error(err), zap.String("name", plan.GetName()))
			}
		}

		if plan.NotifyWithMsg != nil && msg != "" {
			err = plan.NotifyWithMsg(msg)
			if err != nil {
				zlog.Error("NotifyWithMsg failure", zap.Error(err), zap.String("name", plan.GetName()))
			}
		}

		plan.SetStatus(domain.Completed)
	})
	if err != nil {
		panic(fmt.Errorf("register plan failure, nest error: %v", err))
	}

	s.plans = append(s.plans, plan)
}

func (s *scheduler) Start() error {
	if _, err := s.cron.AddFunc("00 16 * * MON,TUE,WED,THU,FRI", func() {
		for _, plan := range s.plans {
			plan.Reset()
		}
	}); err != nil {
		return err
	}

	s.cron.Start()
	return nil
}
