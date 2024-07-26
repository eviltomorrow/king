package controller

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	"github.com/eviltomorrow/king/apps/king-storage/domain/db"
	"github.com/eviltomorrow/king/apps/king-storage/domain/service"
	"github.com/eviltomorrow/king/lib/db/mysql"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/zlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Storage struct {
	pb.UnimplementedStorageServer
}

func NewStorage() *Storage {
	return &Storage{}
}

func (g *Storage) Service() func(*grpc.Server) {
	return func(server *grpc.Server) {
		pb.RegisterStorageServer(server, g)
	}
}

const PER_COMMIT_LIMIT = 32

var sema = semaphore.NewWeighted(32)

func (g *Storage) ArchiveMetadata(ps pb.Storage_ArchiveMetadataServer) error {
	type MetadataWrapper struct {
		Date time.Time
		Data []*model.Metadata
	}
	data := make(map[string]MetadataWrapper)

	var (
		as atomic.Int64
		ad atomic.Int64
		aw atomic.Int64

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

				// _, newspan := opentrace.DefaultTracer().Start(ps.Context(), "StoreMetadata")
				// defer newspan.End()
				s, d, w, err := service.ArchiveMetadata(wrapper.Date, ch)
				if err != nil {
					zlog.Error("Store metadata failure", zap.Error(err))
					return
				}
				as.Add(s)
				ad.Add(d)
				aw.Add(w)
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

			// _, newspan := opentrace.DefaultTracer().Start(ps.Context(), "StoreMetadata")
			// defer newspan.End()
			s, d, w, err := service.ArchiveMetadata(wrapper.Date, ch)
			if err != nil {
				zlog.Error("Store metadata failure", zap.Error(err))
				// newspan.RecordError(err)
				return
			}
			as.Add(s)
			ad.Add(d)
			aw.Add(w)
		}()
		for _, md := range wrapper.Data {
			ch <- md
		}
		close(ch)

		wrapper.Data = wrapper.Data[:0]
	}

	wg.Wait()

	return ps.SendAndClose(&pb.ArchiveResponse{
		Affected: &pb.ArchiveResponse_Affected{
			Stock:     as.Load(),
			QuoteDay:  ad.Load(),
			QuoteWeek: aw.Load(),
		},
	})
}

func (g *Storage) GetStockAll(_ *emptypb.Empty, gs pb.Storage_GetStockAllServer) error {
	var (
		offset, limit int64 = 0, 100
		timeout             = 10 * time.Second
	)

	for {
		stocks, err := db.StockWithSelectRange(mysql.DB, offset, limit, timeout)
		if err != nil {
			return status.Error(codes.InvalidArgument, "req is nil")
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

func (g *Storage) GetQuotesLatest(req *pb.GetQuoteLatestRequest, resp pb.Storage_GetQuoteLatestServer) error {
	var (
		limit   int64 = req.Limit
		kind    string
		timeout = 10 * time.Second
	)
	if limit > 250 {
		return fmt.Errorf("limit should be less than 250")
	}

	switch req.Kind {
	case pb.GetQuoteLatestRequest_Day:
		kind = db.Day
	case pb.GetQuoteLatestRequest_Week:
		kind = db.Week
	default:
		kind = db.Day
	}

	quotes, err := db.QuoteWithSelectManyLatest(mysql.DB, kind, req.Code, req.Date, limit, timeout)
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
