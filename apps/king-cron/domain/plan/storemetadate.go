package plan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
			stub, shutdown, err := client.NewCollectorWithEtcd()
			if err != nil {
				return "", err
			}
			defer shutdown()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			now := time.Now()
			resp, err := stub.StoreMetadata(ctx, &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
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
			if _, err := db.SchedulerRecordWithInsertOne(ctx, mysql.DB, record); err != nil {
				return "", err
			}
			zlog.Info("store metadata success", zap.String("scheduler_id", schedulerId), zap.Int64("stock", resp.Affected.Stock), zap.Int64("quote", resp.Affected.Quote))

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
