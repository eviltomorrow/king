package notification

import (
	"context"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
)

func SendEmail(ctx context.Context, name, address, subject, body string) error {
	stub, closeFunc, err := client.NewEmailWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if _, err = stub.Send(ctx, &pb.Mail{
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
	stub, closeFunc, err := client.NewNtfyWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if _, err = stub.Send(ctx, &pb.Msg{
		// Topic:    "SrxOPwCBiRWZUOq0",
		Topic:    topic,
		Message:  msg,
		Title:    title,
		Priority: 4,
		Tags:     tags,
	}); err != nil {
		return err
	}

	return nil
}
