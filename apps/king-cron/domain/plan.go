package domain

import (
	"context"
	"fmt"
	"sync"

	"github.com/eviltomorrow/king/apps/king-cron/domain/notification"
	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
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
	Todo         func(string) error
	WriteToDB    func(string, error) error

	mutex  sync.Mutex
	Alias  string
	Name   string
	Status StatusCode

	NotifyWithError func(error) error
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

func NotifyForNtfy(templateName string, title string, data map[string]string) error {
	resp, err := client.DefaultTemplate.Render(context.Background(), &pb.RenderRequest{
		TemplateName: templateName,
		Data:         data,
	})
	if err != nil {
		return err
	}
	return notification.DefaultNotifyForNtfyWithMsg(title, resp.Value, []string{"简报", "股票", "统计"})
}

func NotifyForEmail(templateName string, title string, data map[string]string) error {
	resp, err := client.DefaultTemplate.Render(context.Background(), &pb.RenderRequest{
		TemplateName: templateName,
		Data:         data,
	})
	if err != nil {
		return err
	}
	return notification.DefaultNotifyForEmailWithMsg(title, resp.Value)
}
