package server

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/service/db"
	"github.com/eviltomorrow/king/apps/king-storage/service/storage"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	"github.com/eviltomorrow/king/lib/db/mysql"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/grpc/server"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPC struct {
	EtcdClient *clientv3.Client
	Host       string
	Port       int
	AppName    string

	helper *server.GrpcHelper

	pb.UnimplementedStorageServer
}

const PER_COMMIT_LIMIT = 32

var sema = semaphore.NewWeighted(32)

func (g *GRPC) PushMetadata(ps pb.Storage_PushMetadataServer) error {
	type MetadataWrapper struct {
		Date time.Time
		Data []*model.Metadata
	}
	data := make(map[string]MetadataWrapper)

	var (
		affectedS atomic.Int64
		affectedD atomic.Int64
		affectedW atomic.Int64

		wg sync.WaitGroup
	)

	for {
		md, err := ps.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		wrapper, ok := data[md.Date]
		if !ok {
			d, err := time.Parse("2006-01-02", md.Date)
			if err != nil {
				return err
			}
			wrapper = MetadataWrapper{
				Date: d,
				Data: make([]*model.Metadata, 0, 32),
			}
			data[md.Date] = wrapper
		}
		wrapper.Data = append(wrapper.Data, &model.Metadata{
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
		})

		if len(wrapper.Data) == PER_COMMIT_LIMIT {
			ch := make(chan *model.Metadata, PER_COMMIT_LIMIT)

			sema.Acquire(context.Background(), 1)
			wg.Add(1)
			go func() {
				defer sema.Release(1)
				defer wg.Done()
				s, d, w, err := storage.StoreMetadata(wrapper.Date, ch)
				if err != nil {
					zlog.Error("Store metadata failure", zap.Error(err))
					return
				}
				affectedS.Add(s)
				affectedD.Add(d)
				affectedW.Add(w)
			}()
			for _, md := range wrapper.Data {
				ch <- md
			}
			close(ch)

			wrapper.Data = wrapper.Data[:0]
		}
		data[md.Date] = wrapper
	}

	for _, wrapper := range data {
		ch := make(chan *model.Metadata, PER_COMMIT_LIMIT)

		sema.Acquire(context.Background(), 1)
		wg.Add(1)
		go func() {
			defer sema.Release(1)
			defer wg.Done()
			s, d, w, err := storage.StoreMetadata(wrapper.Date, ch)
			if err != nil {
				zlog.Error("Store metadata failure", zap.Error(err))
				return
			}
			affectedS.Add(s)
			affectedD.Add(d)
			affectedW.Add(w)
		}()
		for _, md := range wrapper.Data {
			ch <- md
		}
		close(ch)

		wrapper.Data = wrapper.Data[:0]
	}

	wg.Wait()
	return ps.SendAndClose(&pb.Stats{StockAffected: affectedS.Load(), QuoteDayAffected: affectedD.Load(), QuoteWeekAffected: affectedW.Load()})
}

func (g *GRPC) GetStockFull(_ *emptypb.Empty, gs pb.Storage_GetStockFullServer) error {
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

func (g *GRPC) GetQuoteLatest(req *pb.QuoteRequest, resp pb.Storage_GetQuoteLatestServer) error {
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
			pb.RegisterStorageServer(s, g)
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
