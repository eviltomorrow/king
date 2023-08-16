package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/eviltomorrow/king/app/king-repository/service/data"
	"github.com/eviltomorrow/king/app/king-repository/service/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-repository"
	"github.com/eviltomorrow/king/lib/grpc/server"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string

	helper *server.GrpcHelper

	pb.UnimplementedRepositoryServer
}

func (g *GRPC) ArchiveMetadata(ctx context.Context, req *wrapperspb.StringValue) (*pb.Counter, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request, date is nil")
	}
	d, err := time.Parse("2006-01-02", req.Value)
	if err != nil {
		return nil, err
	}

	stub, closeFunc, err := client.NewCollectorWithEtcd()
	if err != nil {
		return nil, err
	}
	defer closeFunc()

	resp, err := stub.FetchMetadata(context.Background(), &wrapperspb.StringValue{Value: req.Value})
	if err != nil {
		return nil, err
	}

	var (
		pipe   = make(chan *model.Metadata, 128)
		signal = make(chan struct{}, 1)
	)

	go func() {
		defer func() {
			close(pipe)
		}()
		for {
			md, err := resp.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				zlog.Error("FetchMetadata recv failure", zap.Error(err))
				return
			}
			select {
			case <-signal:
				return
			default:
				pipe <- &model.Metadata{
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
				}
			}
		}
	}()
	affectedS, affectedD, affectedW, err := data.TransmissionMetadata(d, pipe)
	if err != nil {
		signal <- struct{}{}
		return nil, err
	}
	return &pb.Counter{AffectedStock: affectedS, AffectedQuoteDay: affectedD, AffectedQuoteWeek: affectedW}, nil
}

func (g *GRPC) GetStockFull(_ *emptypb.Empty, gs pb.Repository_GetStockFullServer) error {
	var (
		offset, limit int64 = 0, 100
		timeout             = 10 * time.Second
	)

	for {
		stocks, err := db.StockWithSelectRange(mysql.DB, offset, limit, timeout)
		if err != nil {
			return err
		}

		for _, stock := range stocks {
			if err := gs.Send(&pb.Stock{Code: stock.Code, Name: stock.Name, Suspend: stock.Suspend}); err != nil {
				return err
			}
		}

		if int64(len(stocks)) < limit {
			break
		}
		offset += limit
	}
	return nil
}

func (g *GRPC) GetQuoteLatest(req *pb.QuoteRequest, resp pb.Repository_GetQuoteLatestServer) error {
	var (
		limit   int64 = req.Limit
		mode    string
		timeout = 10 * time.Second
	)
	if limit > 250 {
		return fmt.Errorf("limit should be less than 250")
	}

	switch req.Mode {
	case pb.QuoteRequest_Day:
		mode = db.Day
	case pb.QuoteRequest_Week:
		mode = db.Week
	default:
		mode = db.Day
	}

	quotes, err := db.QuoteWithSelectManyLatest(mysql.DB, mode, req.Code, req.Date, limit, timeout)
	if err != nil {
		return err
	}

	for _, quote := range quotes {
		if err := resp.Send(&pb.Quote{
			Code:            quote.Code,
			Open:            quote.Open,
			Close:           quote.Close,
			High:            quote.High,
			Low:             quote.Low,
			YesterdayClosed: quote.YesterdayClosed,
			Volume:          quote.Volume,
			Account:         quote.Account,
			Date:            quote.Date.Format("2006-01-02"),
			NumOfYear:       int32(quote.NumOfYear),
		}); err != nil {
			return err
		}
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
			pb.RegisterRepositoryServer(s, g)
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
