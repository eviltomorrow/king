package impl

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-notification/conf"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/smtp"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type EmailServer struct {
	SMTP *conf.SMTP
	pb.UnimplementedEmailServer
}

func (s *EmailServer) Send(ctx context.Context, req *pb.Mail) (*wrapperspb.StringValue, error) {
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
			Name:    s.SMTP.Alias,
			Address: s.SMTP.Username,
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

	if err := smtp.SendWithSSL(s.SMTP.Server, s.SMTP.Username, s.SMTP.Password, message); err != nil {
		return nil, err
	}

	uid := uuid.New()
	zlog.Info("Send email success", zap.String("id", uid.String()), zap.String("msg", message.String()))
	return &wrapperspb.StringValue{Value: uid.String()}, nil
}
