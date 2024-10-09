package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain/notification"
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
	p.SetStatus(Pending)
}

func (p *Plan) IsCompleted() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Status == Completed
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
	Pending StatusCode = iota
	Ready
	Completed
)

func DefaultNotifyWithError(title string, err error, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var e error
	if err := notification.SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, err.Error()); err != nil {
		e = errors.Join(e, fmt.Errorf("send email faiulure, nest error: %v", err))
	}
	if err := notification.SendNtfy(ctx, title, err.Error(), "SrxOPwCBiRWZUOq0", tags); err != nil {
		e = errors.Join(e, fmt.Errorf("send ntfy faiulure, nest error: %v", err))
	}
	return e
}

func DefaultNotifyWithMsg(title, body string, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var e error
	if err := notification.SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, body); err != nil {
		e = errors.Join(e, fmt.Errorf("send email faiulure, nest error: %v", err))
	}
	if err := notification.SendNtfy(ctx, title, body, "SrxOPwCBiRWZUOq0", tags); err != nil {
		e = errors.Join(e, fmt.Errorf("send ntfy faiulure, nest error: %v", err))
	}
	return e
}
