package controller

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-notification/conf"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/ntfy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Ntfy struct {
	NTFY *conf.NTFY
	pb.UnimplementedNtfyServer
}

func NewNtfy(ntfy *conf.NTFY) *Ntfy {
	return &Ntfy{NTFY: ntfy}
}

func (n *Ntfy) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterNtfyServer(server, n)
	}
}

func (n *Ntfy) Send(ctx context.Context, msg *pb.Msg) (*wrapperspb.StringValue, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "msg is nil")
	}
	if msg.Topic == "" {
		return nil, status.Error(codes.InvalidArgument, "topic is nil")
	}

	data, err := ntfy.Send(fmt.Sprintf("%s://%s:%d", n.NTFY.Scheme, n.NTFY.Server, n.NTFY.Port), n.NTFY.Username, n.NTFY.Password, msg.Topic, &ntfy.Msg{
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
