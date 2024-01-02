package impl

import (
	"context"

	"github.com/eviltomorrow/king/apps/king-notification/conf"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/ntfy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type NotificationServer struct {
	NFTY *conf.NFTY
	pb.UnimplementedNotificationServer
}

func (s *NotificationServer) Send(ctx context.Context, msg *pb.Msg) (*wrapperspb.StringValue, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "msg is nil")
	}
	if msg.Topic == "" {
		return nil, status.Error(codes.InvalidArgument, "topic is nil")
	}

	data, err := ntfy.Send(s.NFTY.Server, s.NFTY.Username, s.NFTY.Password, msg.Topic, &ntfy.Msg{
		Title:    msg.Title,
		Message:  msg.Message,
		Priority: msg.Priority,
		Tags:     msg.Tags,
		Attach:   msg.Attach,
	})
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: data}, nil
}
