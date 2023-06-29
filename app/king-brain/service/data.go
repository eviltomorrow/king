package service

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/eviltomorrow/king/app/king-brain/service/calculate"
	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-repository"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DateType int

const (
	_ DateType = iota
	DAY
	WEEK
)

type DataWrapper struct {
	Stock *pb.Stock `json:"stock"`
	Data  []*Data   `json:"data"`
	Type  DateType  `json:"type"`
}

func (d *DataWrapper) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(d)
	if err != nil {
		return fmt.Sprintf("marshal failure, nest error: %v", err)
	}
	return string(buf)
}

type Data struct {
	Quote *pb.Quote `json:"quote"`

	MA_10 float64 `json:"ma_10,omitempty"`
	MA_50 float64 `json:"ma_50,omitempty"`
	MA150 float64 `json:"ma150,omitempty"`
	MA200 float64 `json:"ma200,omitempty"`
}

func (d *Data) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(d)
	if err != nil {
		return fmt.Sprintf("marshal failure, nest error: %v", err)
	}
	return string(buf)
}

func NewDataWrapperChannel(mode DateType, date string) (chan *DataWrapper, error) {
	var data = make(chan *DataWrapper, 64)

	stub, closeFunc, err := client.NewRepository()
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
			wg sync.WaitGroup

			limiter = func() chan struct{} {
				var (
					size = 16
					ch   = make(chan struct{}, size)
				)
				for i := 0; i < size; i++ {
					ch <- struct{}{}
				}
				return ch
			}()
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

			<-limiter
			wg.Add(1)
			go func(stock *pb.Stock) {
				var (
					quotes []*pb.Quote
					err    error
				)
				switch mode {
				case DAY:
					quotes, err = getQuoteDays(stub, stock.Code, date, limit)
					if err != nil {
						zlog.Error("GetQuoteDays failure, nest error: Recv from channel", zap.Error(err))
					} else {
						calculate.ReverseQuote(quotes)
						data <- buildDataWrapper(DAY, stock, quotes)
					}

				case WEEK:
					quotes, err = getQuoteWeeks(stub, stock.Code, date, limit)
					if err != nil {
						zlog.Error("getQuoteWeeks failure, nest error: Recv from channel", zap.Error(err))
					} else {
						calculate.ReverseQuote(quotes)
						data <- buildDataWrapper(WEEK, stock, quotes)
					}

				default:
				}

				wg.Done()
				limiter <- struct{}{}
			}(stock)
		}
		wg.Wait()
		closeFunc()
		close(data)
	}()
	return data, nil
}

func getQuoteDays(stub pb.RepositoryClient, code string, today string, limit int64) ([]*pb.Quote, error) {
	respQuote, err := stub.GetQuoteLatest(context.Background(), &pb.QuoteRequest{
		Code:  code,
		Date:  today,
		Limit: limit,
		Mode:  pb.QuoteRequest_Day,
	})
	if err != nil {
		return nil, err
	}

	var days = make([]*pb.Quote, 0, limit)
	for {
		quote, err := respQuote.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		days = append(days, quote)
	}
	return days, nil
}

func getQuoteWeeks(stub pb.RepositoryClient, code string, today string, limit int64) ([]*pb.Quote, error) {
	respQuote, err := stub.GetQuoteLatest(context.Background(), &pb.QuoteRequest{
		Code:  code,
		Date:  today,
		Limit: limit,
		Mode:  pb.QuoteRequest_Week,
	})
	if err != nil {
		return nil, err
	}

	var days = make([]*pb.Quote, 0, limit)
	for {
		quote, err := respQuote.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		days = append(days, quote)
	}
	return days, nil
}

func buildDataWrapper(mode DateType, stock *pb.Stock, quotes []*pb.Quote) *DataWrapper {
	var (
		closed = make([]float64, 0, len(quotes))
		data   = make([]*Data, 0, len(quotes))
	)

	for _, quote := range quotes {
		closed = append(closed, quote.Close)
		var d = &Data{
			Quote: quote,
			MA_10: maN(closed, 10),
			MA_50: maN(closed, 50),
			MA150: maN(closed, 150),
			MA200: maN(closed, 200),
		}
		data = append(data, d)
	}
	var dw = &DataWrapper{
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
	return mathutil.Trunc2(calculate.MA(closed[len(closed)-n:]))
}
