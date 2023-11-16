package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/netutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AppOption interface {
	apply(*appOptions)
}

type appOptions struct {
	accessHost, listenHost string
	port                   int
	etcdClient             *clientv3.Client
	appName                string
	hookServerFunc         func(*grpc.Server)
}

var defaultAppOptions = appOptions{}

type funcAppOption struct {
	f func(*appOptions)
}

func (fao *funcAppOption) apply(ao *appOptions) {
	fao.f(ao)
}

func WithAccessHost(host string) AppOption {
	return &funcAppOption{
		f: func(ao *appOptions) {
			ao.accessHost = host
		},
	}
}

func WithListenHost(host string) AppOption {
	return &funcAppOption{
		f: func(ao *appOptions) {
			ao.listenHost = host
		},
	}
}

func WithPort(port int) AppOption {
	return &funcAppOption{
		f: func(ao *appOptions) {
			ao.port = port
		},
	}
}

func WithEtcdClient(client *clientv3.Client) AppOption {
	return &funcAppOption{
		f: func(ao *appOptions) {
			ao.etcdClient = client
		},
	}
}

func WithAppName(name string) AppOption {
	return &funcAppOption{
		f: func(ao *appOptions) {
			ao.appName = name
		},
	}
}

func WithRegisterServerFunc(f func(*grpc.Server)) AppOption {
	return &funcAppOption{
		f: func(ao *appOptions) {
			ao.hookServerFunc = f
		},
	}
}

type GrpcHelper struct {
	appOptions *appOptions

	ctx    context.Context
	cancel func()

	revokeFunc func() error
	server     *grpc.Server
}

func NewGrpcHelper(opt ...AppOption) *GrpcHelper {
	opts := &defaultAppOptions
	for _, o := range opt {
		o.apply(opts)
	}

	return &GrpcHelper{appOptions: opts}
}

func (g *GrpcHelper) Init() error {
	var (
		listenHost, accessHost = g.appOptions.listenHost, g.appOptions.accessHost
		port                   = g.appOptions.port
	)

	if listenHost == "" || listenHost == "0.0.0.0" {
		localIp, err := netutil.GetLocalIP2()
		if err != nil {
			return err
		}
		listenHost = localIp
		accessHost = listenHost
	}
	if port == 0 {
		randomPort, err := netutil.GetAvailablePort()
		if err != nil {
			return err
		}
		port = randomPort
	}
	if accessHost == "" {
		accessHost = listenHost
	}

	if g.appOptions.appName == "" || g.appOptions.hookServerFunc == nil {
		return fmt.Errorf("panic: registerAppName or registerServerFunc is nil")
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenHost, port))
	if err != nil {
		return err
	}
	g.server = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			middleware.UnaryServerRecoveryInterceptor,
			middleware.UnaryServerLogInterceptor,
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamServerRecoveryInterceptor,
			middleware.StreamServerLogInterceptor,
		),
	)

	g.ctx, g.cancel = context.WithCancel(context.Background())
	if g.appOptions.etcdClient != nil {
		g.revokeFunc, err = etcd.RegisterService(g.ctx, g.appOptions.appName, accessHost, port, 10, g.appOptions.etcdClient)
		if err != nil {
			return err
		}
	}

	reflection.Register(g.server)
	if g.appOptions.hookServerFunc != nil {
		g.appOptions.hookServerFunc(g.server)
	}

	go func() {
		if err := g.server.Serve(listen); err != nil {
			log.Fatalf("Start server(grpc) failure, nest error: %v", err)
		}
	}()
	return nil
}

func (g *GrpcHelper) Stop() error {
	if g.revokeFunc != nil {
		g.revokeFunc()
	}

	if g.server != nil {
		g.server.Stop()
	}
	g.cancel()
	return nil
}
