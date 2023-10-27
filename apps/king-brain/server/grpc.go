package server

import (
	"context"

	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"github.com/eviltomorrow/king/lib/grpc/server"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string

	helper *server.GrpcHelper

	pb.UnimplementedFinderServer
}

// DiscoverPossibleChance(context.Context, *wrapperspb.StringValue) (*PossibleChance, error)
// CreateBuyPlan(context.Context, *PossibleChance) (*BuyPlan, error)
// FollowReturnWithBuyPlan(context.Context, *BuyPlan) (*PositionPlan, error)

func (g *GRPC) DiscoverPossibleChance(req *wrapperspb.StringValue, ds pb.Finder_DiscoverPossibleChanceServer) error {
	return nil
}

func (g *GRPC) CreateBuyPlan(ctx context.Context, req *pb.PossibleChance) (*pb.BuyPlan, error) {
	return nil, nil
}

func (g *GRPC) FollowReturnWithBuyPlan(ctx context.Context, req *pb.BuyPlan) (*pb.PositionPlan, error) {
	return nil, nil
}

func (g *GRPC) Startup() error {
	g.helper = server.NewGrpcHelper(
		server.WithListenHost(g.Host),
		server.WithPort(g.Port),
		server.WithAppName(g.AppName),
		server.WithEtcdClient(g.EtcdClient),
		server.WithRegisterServerFunc(func(s *grpc.Server) {
			pb.RegisterFinderServer(s, g)
		}),
	)
	return g.helper.Init()
}

func (g *GRPC) Stop() error {
	if g.helper != nil {
		return g.helper.Stop()
	}
	return nil
}
