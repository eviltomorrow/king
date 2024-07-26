package controller

import (
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"google.golang.org/grpc"
)

type Assets struct {
	pb.UnimplementedAssetsServer
}

func NewAssets() *Assets {
	return &Assets{}
}

func (g *Assets) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterAssetsServer(server, g)
	}
}
