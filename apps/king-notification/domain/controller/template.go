package controller

import (
	"context"

	"github.com/eviltomorrow/king/apps/king-notification/domain/service"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Template struct {
	pb.UnimplementedTemplateServer
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterTemplateServer(server, t)
	}
}

func (t *Template) Render(ctx context.Context, req *pb.RenderRequest) (*wrapperspb.StringValue, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "req is nil")
	}
	if len(req.TemplateName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "template name is nil")
	}

	data := make(map[string]interface{}, len(req.Data))
	for k, v := range req.Data {
		data[k] = v
	}
	text, err := service.GenerateMarketStatusHTMLText(req.TemplateName, data)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: text}, nil
}
