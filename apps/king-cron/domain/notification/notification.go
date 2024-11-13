package notification

import (
	"context"
	"errors"
	"fmt"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/setting"
)

func DefaultNotifyForEmailWithError(title string, err error) error {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	var e error
	if err := SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, err.Error()); err != nil {
		e = errors.Join(e, fmt.Errorf("send email failure, nest error: %v", err))
	}
	return e
}

func DefaultNotifyForNtfyWithError(title string, err error, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	var e error
	if err := SendNtfy(ctx, title, err.Error(), "SrxOPwCBiRWZUOq0", tags); err != nil {
		e = errors.Join(e, fmt.Errorf("send ntfy failure, nest error: %v", err))
	}
	return e
}

func DefaultNotifyForEmailWithMsg(title, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	var e error
	if err := SendEmail(ctx, "shepard", "eviltomorrow@163.com", title, body); err != nil {
		e = errors.Join(e, fmt.Errorf("send email failure, nest error: %v", err))
	}
	return e
}

func DefaultNotifyForNtfyWithMsg(title, body string, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	var e error
	if err := SendNtfy(ctx, title, body, "SrxOPwCBiRWZUOq0", tags); err != nil {
		e = errors.Join(e, fmt.Errorf("send ntfy failure, nest error: %v", err))
	}
	return e
}

func SendEmail(ctx context.Context, name, address, subject, body string) error {
	if _, err := client.DefaultEmail.Send(ctx, &pb.Mail{
		To: []*pb.Contact{
			{Name: name, Address: address},
		},
		Subject: subject,
		Body:    body,
	}); err != nil {
		return err
	}

	return nil
}

func SendNtfy(ctx context.Context, title, msg string, topic string, tags []string) error {
	if _, err := client.DefaultNTFY.Send(ctx, &pb.Msg{
		// Topic:    "SrxOPwCBiRWZUOq0",
		Topic:    topic,
		Message:  msg,
		Title:    title,
		Priority: 3,
		Tags:     tags,
	}); err != nil {
		return err
	}

	return nil
}
