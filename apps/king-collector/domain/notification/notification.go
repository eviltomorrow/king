package notification

import (
	"context"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
)

func SendEmail(ctx context.Context, subject, reason string) error {
	stub, closeFunc, err := client.NewEmailWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if _, err = stub.Send(ctx, &pb.Mail{
		To: []*pb.Contact{
			{Name: "Shepard", Address: "eviltomorrow@163.com"},
		},
		Subject: subject,
		Body:    reason,
	}); err != nil {
		return err
	}

	return nil
}

func SendNtfy(ctx context.Context, title, reason string) error {
	stub, closeFunc, err := client.NewNtfyWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if _, err = stub.Send(ctx, &pb.Msg{
		Topic:    "SrxOPwCBiRWZUOq0",
		Message:  reason,
		Title:    title,
		Priority: 4,
		Tags:     []string{"warning", "metadata", "crawl"},
	}); err != nil {
		return err
	}

	return nil
}
