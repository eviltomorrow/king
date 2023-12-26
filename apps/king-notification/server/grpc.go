package server

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-notification/conf"
	"github.com/eviltomorrow/king/apps/king-notification/service"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/grpc/server"
	"github.com/eviltomorrow/king/lib/smtp"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/flosch/pongo2/v6"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string
	SMTP       *conf.SMTP

	helper *server.GrpcHelper

	pb.UnimplementedEmailServer
}

func (g *GRPC) Send(ctx context.Context, req *pb.Mail) (*wrapperspb.StringValue, error) {
	if req == nil {
		return nil, fmt.Errorf("illegal parameter, nest error: mail is nil")
	}
	if len(req.To) == 0 {
		return nil, fmt.Errorf("illegal parameter, nest error: to is nil")
	}

	contentType := smtp.TextHTML
	switch req.ContentType {
	case pb.Mail_TEXT_PLAIN:
		contentType = smtp.TextPlain
	default:
	}
	message := &smtp.Message{
		From: smtp.Contact{
			Name:    g.SMTP.Alias,
			Address: g.SMTP.Username,
		},
		Subject:     req.Subject,
		Body:        req.Body,
		ContentType: contentType,
	}

	to := make([]smtp.Contact, 0, len(req.To))
	for _, c := range req.To {
		if c != nil {
			to = append(to, smtp.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.To = to

	cc := make([]smtp.Contact, 0, len(req.Cc))
	for _, c := range req.Cc {
		if c != nil {
			cc = append(cc, smtp.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.Cc = cc

	bcc := make([]smtp.Contact, 0, len(req.Bcc))
	for _, c := range req.Bcc {
		if c != nil {
			bcc = append(bcc, smtp.Contact{Name: c.Name, Address: c.Address})
		}
	}
	message.Bcc = bcc

	if err := smtp.SendWithSSL(g.SMTP.Server, g.SMTP.Username, g.SMTP.Password, message); err != nil {
		return nil, err
	}

	uid := uuid.New()
	zlog.Info("Send email success", zap.String("id", uid.String()), zap.String("msg", message.String()))
	return &wrapperspb.StringValue{Value: uid.String()}, nil
}

func (g *GRPC) RenderTemplate(ctx context.Context, req *pb.TemplateData) (*wrapperspb.StringValue, error) {
	if req == nil {
		return nil, fmt.Errorf("illegal parameter, nest error: template_data is nil")
	}

	tpl, ok := service.GetTemplate(req.Name)
	if !ok {
		return nil, fmt.Errorf("not found template, nest name: %v", req.Name)
	}

	data := func() pongo2.Context {
		context := make(pongo2.Context, 4)
		for k, v := range req.Data {
			context[k] = v
		}
		return context
	}()
	value, err := tpl.Execute(data)
	return &wrapperspb.StringValue{Value: value}, err
}

func (g *GRPC) Startup() error {
	g.helper = server.NewGrpcHelper(
		server.WithListenHost(g.Host),
		server.WithPort(g.Port),
		server.WithAppName(g.AppName),
		server.WithEtcdClient(g.EtcdClient),
		server.WithRegisterServerFunc(func(s *grpc.Server) {
			pb.RegisterEmailServer(s, g)
		}),
	)
	return g.helper.Init()
}

func (g *GRPC) Stop() error {
	if g.helper != nil {
		return g.helper.Stop()
	}
	return nil
}
