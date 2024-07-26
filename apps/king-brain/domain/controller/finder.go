package controller

import (
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"google.golang.org/grpc"
)

type Finder struct {
	pb.UnimplementedFinderServer
}

func NewFinder() *Finder {
	return &Finder{}
}

func (c *Finder) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterFinderServer(server, c)
	}
}
