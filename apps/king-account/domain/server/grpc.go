package server

import (
	"github.com/eviltomorrow/king/apps/king-account/domain/server/impl"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-account"
	"github.com/eviltomorrow/king/lib/grpc/server"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string

	AssetsServer *impl.AssetsServer

	helper *server.GrpcHelper
}

func (g *GRPC) Startup() error {
	g.helper = server.NewGrpcHelper(
		server.WithListenHost(g.Host),
		server.WithPort(g.Port),
		server.WithAppName(g.AppName),
		server.WithEtcdClient(g.EtcdClient),
		server.WithRegisterServerFunc(func(s *grpc.Server) {
			pb.RegisterAssetsServer(s, g.AssetsServer)
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
