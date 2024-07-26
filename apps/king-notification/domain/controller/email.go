package controller

import (
	"context"

	"github.com/eviltomorrow/king/apps/king-notification/conf"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/smtp"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Email struct {
	SMTP *conf.SMTP

	pb.UnimplementedEmailServer
}

func NewEmail(smtp *conf.SMTP) *Email {
	return &Email{SMTP: smtp}
}

func (e *Email) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterEmailServer(server, e)
	}
}

func (e *Email) Send(ctx context.Context, req *pb.Mail) (*wrapperspb.StringValue, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "mail is nil")
	}
	if len(req.To) == 0 {
		return nil, status.Error(codes.InvalidArgument, "tobody is nil")
	}

	contentType := smtp.TextHTML
	switch req.ContentType {
	case pb.Mail_TEXT_PLAIN:
		contentType = smtp.TextPlain
	default:
	}
	message := &smtp.Message{
		From: smtp.Contact{
			Name:    e.SMTP.Alias,
			Address: e.SMTP.Username,
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

	if err := smtp.SendWithSSL(e.SMTP.Server, e.SMTP.Username, e.SMTP.Password, message); err != nil {
		return nil, err
	}

	uid := uuid.New()
	return &wrapperspb.StringValue{Value: uid.String()}, nil
}
