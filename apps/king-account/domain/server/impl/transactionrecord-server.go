package impl

import (
	"context"
	"fmt"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TransactionRecordServer struct {
	pb.UnimplementedTransactionRecordServer
}

// ListByUserId(context.Context, *wrapperspb.StringValue) (*RecordResp, error)
// ListByOpenID(context.Context, *wrapperspb.StringValue) (*RecordResp, error)
func (s *TransactionRecordServer) ListByUserId(ctx context.Context, req *wrapperspb.StringValue) (*pb.RecordResp, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}
	return nil, nil
}

func (s *TransactionRecordServer) ListByOpenID(ctx context.Context, req *wrapperspb.StringValue) (*pb.RecordResp, error) {
	if req == nil {
		return nil, fmt.Errorf("req is nil")
	}

	return nil, nil
}
