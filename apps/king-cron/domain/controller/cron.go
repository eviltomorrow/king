package controller

import (
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
