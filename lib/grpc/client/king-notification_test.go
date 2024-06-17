package client

import (
	"context"
	"testing"
	"time"

	pb_notification "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
)

func TestSendNotification(t *testing.T) {
	client, closeFunc, err := NewNtfyWithTarget("127.0.0.1:5273")
	if err != nil {
		t.Fatal(err)
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err = client.Send(ctx, &pb_notification.Msg{
		Topic:    "topic_alert",
		Message:  "Hello world",
		Title:    "Hi",
		Priority: 4,
		Tags:     []string{"warning", "metadata", "crawl"},
	}); err != nil {
		t.Fatal(err)
	}
}
