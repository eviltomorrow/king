package controller

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/eviltomorrow/king/apps/king-storage/domain/db"
	"github.com/eviltomorrow/king/apps/king-storage/domain/service"
	"github.com/eviltomorrow/king/lib/db/mysql"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/model"
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

func (g *Storage) PushMetadata(req grpc.ClientStreamingServer[pb.Metadata, pb.PushResponse]) error {
	type MetadataWrapper struct {
		Date time.Time
		Data []*model.Metadata
	}
	data := make(map[string]MetadataWrapper)

	var (
		as atomic.Int64
		aq atomic.Int64
	)

	for {
		md, err := req.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		wrapper, ok := data[md.Date]
		if !ok {
			d, err := time.Parse(time.DateOnly, md.Date)
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
			s, q, err := service.StoreMetadata(wrapper.Date, wrapper.Data)
			if err != nil {
				return err
			}
			as.Add(s)
			aq.Add(q)
			wrapper.Data = wrapper.Data[:0]
		}
		data[md.Date] = wrapper
	}

	for _, wrapper := range data {
		s, q, err := service.StoreMetadata(wrapper.Date, wrapper.Data)
		if err != nil {
			return err
		}
		as.Add(s)
		aq.Add(q)
		wrapper.Data = wrapper.Data[:0]
	}

	return req.SendAndClose(&pb.PushResponse{Affected: &pb.PushResponse_AffectedCount{Stock: as.Load(), Quote: aq.Load()}})
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

func (g *Storage) GetQuoteLatest(req *pb.GetQuoteLatestRequest, resp pb.Storage_GetQuoteLatestServer) error {
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
			Date:            quote.Date.Format(time.DateOnly),
			NumOfYear:       int32(quote.NumOfYear),
		}); err != nil {
			return err
		}
	}
	return nil
}
