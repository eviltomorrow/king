package server

import (
	"github.com/eviltomorrow/king/apps/king-notification/domain/server/impl"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/grpc/server"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string

	EmailServer *impl.EmailServer
	NtfyServer  *impl.NtfyServer

	bootstrap *server.Bootstrap
}

func (g *GRPC) Startup() error {
	g.bootstrap = server.NewGrpcBootstrap(
		server.WithListenHost(g.Host),
		server.WithPort(g.Port),
		server.WithAppName(g.AppName),
		server.WithEtcdClient(g.EtcdClient),
		server.WithRegisterServerFunc(func(s *grpc.Server) {
			pb.RegisterEmailServer(s, g.EmailServer)
			pb.RegisterNtfyServer(s, g.NtfyServer)
		}),
	)
	return g.bootstrap.Init()
}

func (g *GRPC) Stop() error {
	if g.bootstrap != nil {
		return g.bootstrap.Stop()
	}
	return nil
}
