package controller

import (
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"google.golang.org/grpc"
)

type Fund struct {
	pb.UnimplementedFundServer
}

func NewFund() *Fund {
	return &Fund{}
}

func (g *Fund) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterFundServer(server, g)
	}
}
