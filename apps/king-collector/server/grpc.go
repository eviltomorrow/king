package server

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/service/db"
	"github.com/eviltomorrow/king/apps/king-collector/service/synchronize"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
	"github.com/eviltomorrow/king/lib/grpc/server"
	"github.com/eviltomorrow/king/lib/opentrace"
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

	pb.UnimplementedCollectorServer
}

func (g *GRPC) CrawlMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.Counter, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, source is nil")
	}
	if req.Value != "sina" && req.Value != "net126" {
		return nil, fmt.Errorf("invalid request, source is %s", req.Value)
	}
	total, ignore, err := synchronize.DataQuick(req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.Counter{Total: total, Ignore: ignore}, nil
}

func (g *GRPC) FetchMetadata(req *wrapperspb.StringValue, fs pb.Collector_FetchMetadataServer) error {
	if req == nil {
		return fmt.Errorf("invalid request, date is nil")
	}

	d, err := time.Parse(time.DateOnly, req.Value)
	if err != nil {
		return err
	}
	var (
		offset, limit int64 = 0, 100
		lastID        string
		timeout       = 20 * time.Second
	)

	_, span := opentrace.DefaultTracer().Start(fs.Context(), "Loop-SelectMetadataRange")
	defer span.End()

	for {

		metadata, err := db.SelectMetadataRange(mongodb.DB, offset, limit, d.Format(time.DateOnly), lastID, timeout)
		if err != nil {
			span.RecordError(err)
			return err
		}
		for _, md := range metadata {
			if err := fs.Send(&pb.Metadata{
				Source:          md.Source,
				Code:            md.Code,
				Name:            md.Name,
				Open:            md.Open,
				YesterdayClosed: md.YesterdayClosed,
				Latest:          md.Latest,
				High:            md.High,
				Low:             md.Low,
				Volume:          md.Volume,
				Account:         md.Account,
				Date:            md.Date,
				Time:            md.Time,
				Suspend:         md.Suspend,
			}); err != nil {
				span.RecordError(err)
				return err
			}
		}
		if len(metadata) < int(limit) {
			break
		}
		offset += limit
	}

	return nil
}

func (g *GRPC) Startup() error {
	g.helper = server.NewGrpcHelper(
		server.WithListenHost(g.Host),
		server.WithPort(g.Port),
		server.WithAppName(g.AppName),
		server.WithEtcdClient(g.EtcdClient),
		server.WithRegisterServerFunc(func(s *grpc.Server) {
			pb.RegisterCollectorServer(s, g)
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
