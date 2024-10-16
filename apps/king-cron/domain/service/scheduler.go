package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/setting"
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
		now := time.Now()
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
				zlog.Error("precondition check failure", zap.Error(err), zap.Any("status", status), zap.String("name", plan.GetName()))
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
		msg, err = plan.Todo(schedulerId)
		if err != nil {
			zlog.Error("plan execute failure", zap.Error(err), zap.String("name", plan.GetName()))
		} else {
			zlog.Info("plan execute success", zap.String("name", plan.GetName()))
		}

		if plan.NotifyWithError != nil && err != nil {
			err = plan.NotifyWithError(err)
			if err != nil {
				zlog.Error("notify with error failure", zap.Error(err), zap.String("name", plan.GetName()))
			}
		}

		if plan.NotifyWithMsg != nil && msg != "" {
			err = plan.NotifyWithMsg(msg)
			if err != nil {
				zlog.Error("notify with msg failure", zap.Error(err), zap.String("name", plan.GetName()))
			}
		}
		plan.SetStatus(domain.Completed)

		// 数据库
		currentStatus, currentCode := func() (string, sql.NullString) {
			if plan.Type == domain.ASYNC {
				return domain.StatusProcessing, sql.NullString{}
			}
			if err == nil {
				return domain.StatusCompleted, sql.NullString{String: codes.SUCCESS, Valid: true}
			} else {
				return domain.StatusCompleted, sql.NullString{String: codes.FAILURE, Valid: true}
			}
		}()

		record := &db.SchedulerRecord{
			Id:          schedulerId,
			Name:        plan.Name,
			Date:        now,
			ServiceName: plan.CallFuncInfo.ServiceName,
			FuncName:    plan.CallFuncInfo.FuncName,
			Code:        currentCode,
			Status:      currentStatus,
		}

		ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10TIMEOUT)
		defer cancel()

		if _, err := db.SchedulerRecordWithInsertOne(ctx, mysql.DB, record); err != nil {
			zlog.Error("scheduler record insert failure, nest error: %v", zap.Error(err), zap.String("name", plan.Name))
		}

	})
	if err != nil {
		return fmt.Errorf("plan register failure, nest error: %v", err)
	}

	s.plans = append(s.plans, plan)
	return nil
}

func (s *scheduler) Start() error {
	if _, err := s.cron.AddFunc("25 16 * * MON,TUE,WED,THU,FRI", func() {
		for _, plan := range s.plans {
			plan.Reset()
		}
	}); err != nil {
		return err
	}

	s.cron.Start()
	return nil
}
