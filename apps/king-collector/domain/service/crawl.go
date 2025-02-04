package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/domain/datasource"
	"github.com/eviltomorrow/king/apps/king-collector/domain/db"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

func CrawlMetadataQuick(ctx context.Context, source string, baseCodeList []string) (int64, int64, error) {
	if len(baseCodeList) == 0 {
		return 0, 0, nil
	}

	select {
	case inFlightSem <- struct{}{}:
		defer func() { <-inFlightSem }()
	default:
		return 0, 0, fmt.Errorf("synchronize metadata busy")
	}

	pipe := make(chan *model.Metadata, 128)
	go persistenceMetadata(ctx, source, pipe)

	options := []func(*model.Metadata) bool{
		func(m *model.Metadata) bool {
			if m != nil && m.Volume != 0 {
				return true
			}
			return false
		},
	}
	return fetchMetadata(nil, baseCodeList, pipe, options)
}

var lastSyncCount int64 = -1

func CrawlMetadataSlow(ctx context.Context, source string, baseCodeList []string, randomPeriod []int) (int64, int64, error) {
	if len(baseCodeList) == 0 {
		return 0, 0, nil
	}

	select {
	case inFlightSem <- struct{}{}:
		defer func() { <-inFlightSem }()
	default:
		return 0, 0, fmt.Errorf("synchronize metadata busy")
	}

	pipe := make(chan *model.Metadata, 128)
	go persistenceMetadata(ctx, source, pipe)

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
	total, ignore, err := fetchMetadata(randomPeriod, baseCodeList, pipe, options)

	if lastSyncCount != -1 && (total > lastSyncCount && float64(total-lastSyncCount) > float64(lastSyncCount)*0.1) {
		return total, ignore, fmt.Errorf("possible missing data, total: %v, lastSyncCount: %v", total, lastSyncCount)
	}

	if total != 0 {
		lastSyncCount = total
	}

	return total, ignore, err
}

func fetchMetadata(randomPeriod []int, baseCodeList []string, pipe chan *model.Metadata, options []func(*model.Metadata) bool) (int64, int64, error) {
	defer func() {
		close(pipe)
	}()

	const (
		size       = 30
		limitTimes = 5
	)
	var (
		codeList                      = make([]string, 0, size)
		retrytimes                    = 0
		totalCount, ignoreCount int64 = 0, 0
	)

	for code := range genCode(baseCodeList) {
		codeList = append(codeList, code)
		if len(codeList) == size {
		retry_1:
			data, err := datasource.FetchMetadataFromSina(codeList)
			if err != nil {
				retrytimes++
				if retrytimes == limitTimes {
					return totalCount, ignoreCount, fmt.Errorf("FetchMeatadata failure, nest error: %v, codeList: %v", err, codeList)
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
						zlog.Warn("metadata ignore", zap.String("data", d.String()))
						ignoreCount++
						continue next1
					}
				}

				pipe <- d
				totalCount++
			}

			if len(randomPeriod) == 2 {
				time.Sleep(time.Duration(mathutil.GenRandInt(randomPeriod[0], randomPeriod[1])) * time.Second)
			} else {
				time.Sleep(300 * time.Millisecond)
			}
		}
	}

	if len(codeList) != 0 {
		retrytimes = 0
	retry_2:
		data, err := datasource.FetchMetadataFromSina(codeList)
		if err != nil {
			retrytimes++
			if retrytimes == limitTimes {
				return totalCount, ignoreCount, fmt.Errorf("fetch metadata failure, nest error: %v, codeList: %v", err, codeList)
			} else {
				time.Sleep(3 * time.Minute)
				goto retry_2
			}
		}
	next2:
		for _, d := range data {
			for _, f := range options {
				if ok := f(d); !ok {
					zlog.Warn("data ignore", zap.String("data", d.String()))
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

func persistenceMetadata(_ context.Context, source string, pipe chan *model.Metadata) {
	const (
		size = 30
	)
	dataList := make([]*model.Metadata, 0, size)
	for data := range pipe {
		if _, err := db.DeleteMetadataByDate(context.Background(), mongodb.DB, source, data.Code, data.Date); err != nil {
			zlog.Error("delete metadata failure", zap.Error(err), zap.String("data", data.String()))
		} else {
			dataList = append(dataList, data)
			if len(dataList) == size {
				if _, err := db.InsertMetadataMany(context.Background(), mongodb.DB, source, dataList); err != nil {
					for _, d := range dataList {
						zlog.Error("insert metadata failure", zap.Error(err), zap.String("data", d.String()))
					}
				}
				dataList = dataList[:0]
			}
		}
	}

	cache := make([]*model.Metadata, 0, len(dataList))
	if len(dataList) != 0 {
		for _, data := range dataList {
			if _, err := db.DeleteMetadataByDate(context.Background(), mongodb.DB, source, data.Code, data.Date); err != nil {
				zlog.Error("delete metadata failure", zap.Error(err), zap.String("data", data.String()))
			} else {
				cache = append(cache, data)
			}
		}
	}

	if len(cache) != 0 {
		if _, err := db.InsertMetadataMany(context.Background(), mongodb.DB, source, cache); err != nil {
			for _, d := range cache {
				zlog.Error("insert metadata failure", zap.Error(err), zap.String("data", d.String()))
			}
		}
	}
}
