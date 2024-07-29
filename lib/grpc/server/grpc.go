package server

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/eviltomorrow/king/lib/buildinfo"
	"github.com/eviltomorrow/king/lib/certificate"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/grpc/middleware"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/resolver"
)

type GRPC struct {
	DisableTLS bool

	AccessIP string
	BindIP   string
	BindPort int

	RegisteredAPI []func(*grpc.Server)

	server     *grpc.Server
	ctx        context.Context
	cancel     func()
	revokeFunc func() error
}

func NewGRPC(c *Config, supported ...func(*grpc.Server)) *GRPC {
	return &GRPC{
		DisableTLS: c.DisableTLS,
		AccessIP:   c.AccessIP,
		BindIP:     c.BindIP,
		BindPort:   c.BindPort,

		RegisteredAPI: supported,
	}
}

func (g *GRPC) Serve() error {
	var creds credentials.TransportCredentials
	if !g.DisableTLS {
		ipList := make([]string, 0, 4)
		ipList = append(ipList, system.Network.BindIP)
		ipList = append(ipList, g.BindIP)
		if system.Network.AccessIP == "" {
			ipList = append(ipList, system.Network.AccessIP)
		}

		err := certificate.CreateOrOverrideFile(certificate.BuildDefaultAppInfo(ipList), &certificate.Config{
			CaCertFile:     filepath.Join(system.Directory.UsrDir, "certs/ca.crt"),
			CaKeyFile:      filepath.Join(system.Directory.UsrDir, "certs/ca.key"),
			ClientCertFile: filepath.Join(system.Directory.VarDir, "certs/client.crt"),
			ClientKeyFile:  filepath.Join(system.Directory.VarDir, "certs/client.pem"),
			ServerCertFile: filepath.Join(system.Directory.VarDir, "certs/server.crt"),
			ServerKeyFile:  filepath.Join(system.Directory.VarDir, "certs/server.pem"),
		})
		if err != nil {
			return err
		}

		creds, err = certificate.LoadServerCredentials(&certificate.Config{
			CaCertFile:     filepath.Join(system.Directory.UsrDir, "certs/ca.crt"),
			ServerCertFile: filepath.Join(system.Directory.VarDir, "certs/server.crt"),
			ServerKeyFile:  filepath.Join(system.Directory.VarDir, "certs/server.pem"),
		})
		if err != nil {
			return err
		}
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.BindIP, g.BindPort))
	if err != nil {
		return err
	}

	g.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryServerRecoveryInterceptor,
			middleware.UnaryServerLogInterceptor,
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamServerRecoveryInterceptor,
			middleware.StreamServerLogInterceptor,
		),
		grpc.Creds(creds),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	reflection.Register(g.server)
	for _, register := range g.RegisteredAPI {
		register(g.server)
	}

	go func() {
		if err := g.server.Serve(listen); err != nil {
			zlog.Fatal("GRPC Server startup failure", zap.Error(err))
		}
	}()

	g.ctx, g.cancel = context.WithCancel(context.Background())
	if etcd.Client != nil {
		resolver.Register(lb.NewBuilder(etcd.Client))
		g.revokeFunc, err = etcd.RegisterService(g.ctx, buildinfo.AppName, system.Network.AccessIP, g.BindPort, 10)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GRPC) Stop() error {
	if g.revokeFunc != nil {
		g.revokeFunc()
	}
	if g.server != nil {
		g.server.GracefulStop()
	}
	g.cancel()

	return nil
}
