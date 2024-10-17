package controller

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-cron"
	"google.golang.org/grpc"
)

type Cron struct {
	pb.UnimplementedSchedulerServer
}

func NewCron() *Cron {
	return &Cron{}
}

func (g *Cron) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterSchedulerServer(server, g)
	}
}

func (g *Cron) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.CallbackResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, req is nil")
	}

	code := func() string {
		switch req.Code {
		case pb.CallbackRequest_SUCCESS:
			return codes.SUCCESS
		case pb.CallbackRequest_FAILURE:
			return codes.FAILURE
		default:
			return codes.FAILURE
		}
	}()
	if _, err := db.SchedulerRecordWithUpdateStatus(ctx, mysql.DB, domain.ProgressCompleted, code, req.ErrorMsg, req.SchedulerId); err != nil {
		return nil, err
	}
	return &pb.CallbackResponse{}, nil
}
