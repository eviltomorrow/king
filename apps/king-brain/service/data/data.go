package data

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/eviltomorrow/king/apps/king-brain/service/calculate"
	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DateType int

const (
	_ DateType = iota
	DAY
	WEEK
)

type DataWrapper struct {
	Type DateType `json:"type"`

	Stock *pb.Stock `json:"stock"`
	Data  []*D      `json:"data"`
}

func (d *DataWrapper) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(d)
	if err != nil {
		return fmt.Sprintf("marshal failure, nest error: %v", err)
	}
	return string(buf)
}

type D struct {
	Quote *pb.Quote `json:"quote"`

	MA_10 float64 `json:"ma_10,omitempty"`
	MA_50 float64 `json:"ma_50,omitempty"`
	MA150 float64 `json:"ma150,omitempty"`
	MA200 float64 `json:"ma200,omitempty"`
}

func (d *D) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(d)
	if err != nil {
		return fmt.Sprintf("marshal failure, nest error: %v", err)
	}
	return string(buf)
}

func NewDataWrapperChannel(mode DateType, date string) (chan *DataWrapper, error) {
	data := make(chan *DataWrapper, 64)

	stub, closeFunc, err := client.NewStorageWithEtcd()
	if err != nil {
		return nil, err
	}

	respStock, err := stub.GetStockFull(context.Background(), &emptypb.Empty{})
	if err != nil {
		closeFunc()
		return nil, err
	}
	go func() {
		var (
			sema = semaphore.NewWeighted(16)
			wg   sync.WaitGroup

			limit int64 = 250
		)
		for {
			stock, err := respStock.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				zlog.Error("GetStockFull failure, nest error: Recv from channel", zap.Error(err))
				break
			}

			sema.Acquire(context.Background(), 1)
			wg.Add(1)
			go func(stock *pb.Stock) {
				var (
					quotes []*pb.Quote
					err    error
				)
				switch mode {
				case DAY:
					quotes, err = getQuotes(stub, pb.QuoteRequest_Day, stock.Code, date, limit)
					if err != nil {
						zlog.Error("GetQuoteDays failure, nest error: Recv from channel", zap.Error(err))
					} else {
						calculate.ReverseQuote(quotes)
						data <- buildDataWrapper(DAY, stock, quotes)
					}

				case WEEK:
					quotes, err = getQuotes(stub, pb.QuoteRequest_Week, stock.Code, date, limit)
					if err != nil {
						zlog.Error("getQuoteWeeks failure, nest error: Recv from channel", zap.Error(err))
					} else {
						calculate.ReverseQuote(quotes)
						data <- buildDataWrapper(WEEK, stock, quotes)
					}

				default:
				}

				wg.Done()
				sema.Release(1)
			}(stock)
		}
		wg.Wait()
		closeFunc()
		close(data)
	}()
	return data, nil
}

func getQuotes(stub pb.StorageClient, mode pb.QuoteRequest_Mode, code string, today string, limit int64) ([]*pb.Quote, error) {
	respQuote, err := stub.GetQuoteLatest(context.Background(), &pb.QuoteRequest{
		Code:  code,
		Date:  today,
		Limit: limit,
		Mode:  mode,
	})
	if err != nil {
		return nil, err
	}

	quotes := make([]*pb.Quote, 0, limit)
	for {
		quote, err := respQuote.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func buildDataWrapper(mode DateType, stock *pb.Stock, quotes []*pb.Quote) *DataWrapper {
	var (
		closed = make([]float64, 0, len(quotes))
		data   = make([]*D, 0, len(quotes))
	)

	for _, quote := range quotes {
		closed = append(closed, quote.Close)
		d := &D{
			Quote: quote,
			MA_10: maN(closed, 10),
			MA_50: maN(closed, 50),
			MA150: maN(closed, 150),
			MA200: maN(closed, 200),
		}
		data = append(data, d)
	}
	dw := &DataWrapper{
		Stock: stock,
		Data:  data,
		Type:  mode,
	}
	return dw
}

func maN(closed []float64, n int) float64 {
	if len(closed) < n {
		return 0
	}
	return mathutil.Trunc4(calculate.MA(closed[len(closed)-n:]))
}
