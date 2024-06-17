package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/domain/persistence"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	client_grpc "github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

func SynchronizeMetadataQuick(ctx context.Context, source string) (int64, int64, error) {
	select {
	case inFlightSem <- struct{}{}:
		defer func() { <-inFlightSem }()
	default:
		return 0, 0, fmt.Errorf("synchronize metadata busy")
	}

	pipe := make(chan *model.Metadata, 128)
	go persistenceMetadata(source, pipe)

	options := []func(*model.Metadata) bool{
		func(m *model.Metadata) bool {
			if m != nil && m.Volume != 0 {
				return true
			}
			return false
		},
	}
	return fetchMetadata(source, false, pipe, options)
}

func SynchronizeMetadataSlow(source string) (int64, int64, error) {
	select {
	case inFlightSem <- struct{}{}:
		defer func() { <-inFlightSem }()
	default:
		return 0, 0, fmt.Errorf("synchronize metadata busy")
	}

	pipe := make(chan *model.Metadata, 128)
	go persistenceMetadata(source, pipe)

	today := time.Now().Format(time.DateOnly)
	options := []func(*model.Metadata) bool{
		func(m *model.Metadata) bool {
			if m != nil && m.Date == today {
				return true
			}
			return false
		},
		func(m *model.Metadata) bool {
			if m != nil && m.Volume != 0 {
				return true
			}
			return false
		},
	}
	return fetchMetadata(source, true, pipe, options)
}

func fetchMetadata(source string, slow bool, pipe chan *model.Metadata, options []func(*model.Metadata) bool) (int64, int64, error) {
	defer func() {
		close(pipe)
	}()

	fetch, ok := FetchFactories[source]
	if !ok {
		return 0, 0, fmt.Errorf("not found fetchFunc, source = [%s]", source)
	}

	var (
		retrytimes                    = 0
		totalCount, ignoreCount int64 = 0, 0
		codeList                      = make([]string, 0, Size)
	)

	for code := range genCode() {
		codeList = append(codeList, code)
		if len(codeList) == Size {
		retry_1:
			data, err := fetch(codeList)
			if err != nil {
				retrytimes++
				if retrytimes == Limit {
					return totalCount, ignoreCount, fmt.Errorf("FetchMeatadata failure, nest error: %v, source: [%v], codeList: %v", err, source, codeList)
				} else {
					time.Sleep(3 * time.Minute)
					goto retry_1
				}
			}
			retrytimes = 0
			codeList = codeList[:0]

		next1:
			for _, d := range data {
				for _, f := range options {
					if ok := f(d); !ok {
						zlog.Warn("Ignore metadata", zap.String("data", d.String()))
						ignoreCount++
						continue next1
					}
				}

				pipe <- d
				totalCount++
			}

			if slow {
				time.Sleep(time.Duration(mathutil.GenRandInt(RandomPeriod[0], RandomPeriod[1])) * time.Second)
			} else {
				time.Sleep(300 * time.Millisecond)
			}
		}
	}

	if len(codeList) != 0 {
		retrytimes = 0
	retry_2:
		data, err := fetch(codeList)
		if err != nil {
			retrytimes++
			if retrytimes == Limit {
				return totalCount, ignoreCount, fmt.Errorf("FetchMetadata failure, nest error: %v, source: [%v], codeList: %v", err, source, codeList)
			} else {
				time.Sleep(3 * time.Minute)
				goto retry_2
			}
		}
	next2:
		for _, d := range data {
			for _, f := range options {
				if ok := f(d); !ok {
					zlog.Warn("Ignore data", zap.String("data", d.String()))
					ignoreCount++
					continue next2
				}
			}

			pipe <- d
			totalCount++
		}
	}
	return totalCount, ignoreCount, nil
}

func persistenceMetadata(source string, pipe chan *model.Metadata) {
	dataList := make([]*model.Metadata, 0, Size)
	for data := range pipe {
		if _, err := persistence.DeleteMetadataByDate(mongodb.DB, source, data.Code, data.Date, Timeout); err != nil {
			zlog.Error("DeleteMetadata by date failure", zap.Error(err), zap.String("data", data.String()))
		} else {
			dataList = append(dataList, data)
			if len(dataList) == Size {
				if _, err := persistence.InsertMetadataMany(mongodb.DB, source, dataList, Timeout); err != nil {
					for _, d := range dataList {
						zlog.Error("InsertMetadata many failure", zap.Error(err), zap.String("data", d.String()))
					}
				}
				dataList = dataList[:0]
			}
		}
	}

	cache := make([]*model.Metadata, 0, len(dataList))
	if len(dataList) != 0 {
		for _, data := range dataList {
			if _, err := persistence.DeleteMetadataByDate(mongodb.DB, source, data.Code, data.Date, Timeout); err != nil {
				zlog.Error("DeleteMetadata by date failure", zap.Error(err), zap.String("data", data.String()))
			} else {
				cache = append(cache, data)
			}
		}
	}

	if len(cache) != 0 {
		if _, err := persistence.InsertMetadataMany(mongodb.DB, source, cache, Timeout); err != nil {
			for _, d := range cache {
				zlog.Error("InsertMetadata many failure", zap.Error(err), zap.String("data", d.String()))
			}
		}
	}
}

func PushMetadataToStorage(ctx context.Context, date string) (int64, int64, int64, int64, error) {
	client, closeFunc, err := client_grpc.NewStorageWithEtcd()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer closeFunc()

	stub, err := client.PushMetadata(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var (
		offset, limit, total int64 = 0, 100, 0
		lastID               string
		timeout              = 20 * time.Second
	)

	for {
		metadata, err := persistence.SelectMetadataRange(mongodb.DB, offset, limit, date, lastID, timeout)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		for _, md := range metadata {
			if err := stub.Send(&pb.Metadata{
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
				return 0, 0, 0, 0, err
			}
			total++
		}
		if len(metadata) < int(limit) {
			break
		}
		offset += limit
	}

	resp, err := stub.CloseAndRecv()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return total, resp.StockAffected, resp.QuoteDayAffected, resp.QuoteWeekAffected, nil
}
