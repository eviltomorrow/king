package domain

import (
	"context"
	"sync"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain/notification"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

type Plan struct {
	Precondition func() (StatusCode, error)
	Todo         func() (string, error)

	mutex  sync.Mutex
	Name   string
	Status StatusCode

	NotifyWithError func(error) error
	NotifyWithMsg   func(string) error
}

func (p *Plan) GetName() string {
	return p.Name
}

func (p *Plan) Reset() {
	p.SetStatus(Ready)
}

func (p *Plan) IsReady() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Status == Ready
}

func (p *Plan) SetStatus(status StatusCode) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Status = status
}

const (
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)

type StatusCode int

const (
	Ready StatusCode = iota
	Pending
	Completed
)

func DefaultNotifyWithError(title string, err error, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := notification.SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, err.Error()); err != nil {
		zlog.Error("Send email failure", zap.Error(err))
	}
	if err := notification.SendNtfy(ctx, title, err.Error(), "SrxOPwCBiRWZUOq0", tags); err != nil {
		zlog.Error("Send ntfy failure", zap.Error(err))
	}
	return nil
}

func DefaultNotifyWithMsg(title, body string, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := notification.SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, body); err != nil {
		zlog.Error("Send email failure", zap.Error(err))
	}
	if err := notification.SendNtfy(ctx, title, body, "SrxOPwCBiRWZUOq0", tags); err != nil {
		zlog.Error("Send ntfy failure", zap.Error(err))
	}
	return nil
}
