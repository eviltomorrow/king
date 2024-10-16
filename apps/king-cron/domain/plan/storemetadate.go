package plan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/setting"
	"github.com/eviltomorrow/king/lib/snowflake"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithStoreMetadata = "存储数据"
)

func CronWithStoreMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10TIMEOUT)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithStoreMetadata, time.Now().Format(time.DateOnly))
			if err != nil && err != sql.ErrNoRows {
				return 0, err
			}
			if record != nil && record.Status == domain.StatusCompleted {
				if record.Code.String == codes.SUCCESS {
					return domain.Completed, nil
				}
				return 0, errors.New(record.ErrorMsg.String)
			}

			record, err = db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithCrawlMetadata, time.Now().Format(time.DateOnly))
			if err != nil {
				return 0, err
			}

			if record.Status == domain.StatusCompleted {
				if record.Code.String == codes.SUCCESS {
					return domain.Ready, nil
				}
				return 0, errors.New(record.ErrorMsg.String)
			}
			return domain.Pending, nil
		},

		Todo: func() (string, error) {
			stubStorage, closeFuncStorage, err := client.NewStorageWithEtcd()
			if err != nil {
				return "", err
			}
			defer closeFuncStorage()

			target, err := stubStorage.PushMetadata(context.Background())
			if err != nil {
				return "", err
			}

			stubCollector, closeFuncCollector, err := client.NewCollectorWithEtcd()
			if err != nil {
				return "", err
			}
			defer closeFuncCollector()

			now := time.Now()
			source, err := stubCollector.FetchMetadata(context.Background(), &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
			if err != nil {
				return "", err
			}
			for {
				md, err := source.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return "", err
				}

				if err := target.Send(md); err != nil {
					return "", err
				}
			}

			resp, err := target.CloseAndRecv()
			if err != nil {
				return "", err
			}

			schedulerId := snowflake.GenerateID()
			record := &db.SchedulerRecord{
				Id:          schedulerId,
				Name:        NameWithStoreMetadata,
				Date:        now,
				ServiceName: "collector",
				FuncName:    "StoreMetadata",
				Status:      domain.StatusCompleted,
			}

			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10TIMEOUT)
			defer cancel()

			if _, err := db.SchedulerRecordWithInsertOne(ctx, mysql.DB, record); err != nil {
				return "", err
			}
			zlog.Info("store metadata success", zap.String("scheduler_id", schedulerId), zap.Int64("stocks", resp.Affected.Stocks), zap.Int64("days", resp.Affected.Days), zap.Int64("weeks", resp.Affected.Weeks))

			return "", nil
		},

		NotifyWithError: func(err error) error {
			return domain.DefaultNotifyWithError(NameWithStoreMetadata, fmt.Errorf("failure: %v", err), []string{"缓存数据", "数据库"})
		},

		Status: domain.Ready,
		Name:   NameWithStoreMetadata,
	}
	return p
}
