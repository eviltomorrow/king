package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/eviltomorrow/king/apps/king-cron/domain/notification"
	"github.com/eviltomorrow/king/lib/setting"
)

var cache = make(map[string]func() *Plan, 32)

func RegisterPlan(name string, f func() *Plan) {
	cache[name] = f
}

func GetPlan(name string) (*Plan, bool) {
	f, ok := cache[name]
	if !ok {
		return nil, false
	}
	return f(), true
}

type CallInfo struct {
	ServiceName string
	FuncName    string
}

type Plan struct {
	Precondition func() (StatusCode, error)
	Todo         func(string) (string, error)
	WriteToDB    func(string, error) error

	mutex  sync.Mutex
	Alias  string
	Name   string
	Status StatusCode

	NotifyWithError func(error) error
	NotifyWithData  func(string) error
}

func (p *Plan) Check() error {
	if p.Todo == nil {
		return fmt.Errorf("plan's todo func is nil")
	}
	if p.WriteToDB == nil {
		return fmt.Errorf("plan's write to db is nil")
	}
	if p.Name == "" {
		return fmt.Errorf("plan's name is nil")
	}
	if p.Status != Ready {
		return fmt.Errorf("plan's status is not ready")
	}
	return nil
}

func (p *Plan) GetName() string {
	return p.Name
}

func (p *Plan) GetAlias() string {
	return p.Alias
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

type StatusCode int

const (
	Pending StatusCode = iota
	Ready
	Completed
)

const (
	ProgressProcessing = "processing"
	ProgressCompleted  = "completed"
)

func DefaultNotifyWithError(title string, err error, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	var e error
	if err := notification.SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, err.Error()); err != nil {
		e = errors.Join(e, fmt.Errorf("send email failure, nest error: %v", err))
	}
	if err := notification.SendNtfy(ctx, title, err.Error(), "SrxOPwCBiRWZUOq0", tags); err != nil {
		e = errors.Join(e, fmt.Errorf("send ntfy failure, nest error: %v", err))
	}
	return e
}

func DefaultNotifyWithMsg(title, body string, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	var e error
	if err := notification.SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, body); err != nil {
		e = errors.Join(e, fmt.Errorf("send email failure, nest error: %v", err))
	}
	if err := notification.SendNtfy(ctx, title, body, "SrxOPwCBiRWZUOq0", tags); err != nil {
		e = errors.Join(e, fmt.Errorf("send ntfy failure, nest error: %v", err))
	}
	return e
}
