package service

import (
	"fmt"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/lib/snowflake"
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

func (s *scheduler) Register(cron string, plan *domain.Plan) error {
	if err := plan.Check(); err != nil {
		return fmt.Errorf("plan check failure, nest error: %v", err)
	}
	_, err := s.cron.AddFunc(cron, func() {
		if plan.IsCompleted() {
			return
		}

		var (
			status = domain.Ready

			err         error
			schedulerId = snowflake.GenerateID()
		)

		if plan.Precondition != nil {
			status, err = plan.Precondition()
			if err != nil {
				zlog.Error("precondition check failure", zap.Error(err), zap.Any("status", status), zap.String("name", plan.GetAlias()))
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

		data := ""
		data, err = plan.Todo(schedulerId)
		if err != nil {
			zlog.Error("plan execute failure", zap.Error(err), zap.String("alias", plan.GetAlias()))
		} else {
			zlog.Info("plan execute success", zap.String("alias", plan.GetAlias()))
		}

		if plan.WriteToDB != nil {
			if err := plan.WriteToDB(schedulerId, err); err != nil {
				zlog.Error("WriteToDB failure", zap.String("alias", plan.Alias), zap.String("schedulerId", schedulerId), zap.Error(err))
			}
		}

		if plan.NotifyWithError != nil && err != nil {
			err := plan.NotifyWithError(err)
			if err != nil {
				zlog.Error("notify with error failure", zap.Error(err), zap.String("name", plan.GetAlias()))
			}
		}

		if plan.NotifyWithData != nil && data != "" {
			err := plan.NotifyWithData(data)
			if err != nil {
				zlog.Error("notify with msg failure", zap.Error(err), zap.String("name", plan.GetAlias()))
			}
		}
		plan.SetStatus(domain.Completed)

	})
	if err != nil {
		return fmt.Errorf("plan register failure, nest error: %v", err)
	}

	s.plans = append(s.plans, plan)
	return nil
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
