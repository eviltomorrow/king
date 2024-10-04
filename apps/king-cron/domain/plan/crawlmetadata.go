package plan

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/snowflake"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithCrawlMetadata = "爬取数据"
	NameWithStoreMetadata = "存储数据"
)

func CronWithCrawlMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: nil,
		Todo: func() error {
			stub, shutdown, err := client.NewCollectorWithEtcd()
			if err != nil {
				return err
			}
			defer shutdown()

			schedulerId := snowflake.GenerateID()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			md := metadata.MD{
				"scheduler_id": []string{schedulerId},
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
			if _, err := stub.CrawlMetadataAsync(ctx, &wrapperspb.StringValue{Value: "sina"}); err != nil {
				return err
			}

			t := time.Now().Format(time.DateOnly)
			now, err := time.Parse(time.DateOnly, t)
			if err != nil {
				return err
			}
			record := &db.SchedulerRecord{
				Id:          schedulerId,
				Name:        NameWithCrawlMetadata,
				Date:        now,
				ServiceName: "collector",
				FuncName:    "CrawlMetadataAsync",
				Status:      domain.StatusProcessing,
			}
			if _, err := db.SchedulerRecordWithInsertOne(ctx, mysql.DB, record); err != nil {
				return err
			}

			zlog.Info("CrawlMetadataAsync success", zap.String("scheduler_id", schedulerId))

			return nil
		},
	}
	return p
}

func CronWithStoreMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: func() (bool, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithStoreMetadata, time.Now().Format(time.DateOnly))
			if err != nil && err != sql.ErrNoRows {
				return false, err
			}
			if record != nil && record.Status == domain.StatusCompleted {
				if record.Code == codes.SUCCESS {
					return false, nil
				}
				return false, errors.New(record.ErrorMsg)
			}

			record, err = db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithCrawlMetadata, time.Now().Format(time.DateOnly))
			if err != nil {
				return false, err
			}

			if record.Status == domain.StatusCompleted {
				if record.Code == codes.SUCCESS {
					return true, nil
				}
				return false, errors.New(record.ErrorMsg)
			}
			return false, nil
		},

		Todo: func() error {
			stub, shutdown, err := client.NewCollectorWithEtcd()
			if err != nil {
				return err
			}
			defer shutdown()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			now := time.Now()
			resp, err := stub.StoreMetadata(ctx, &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
			if err != nil {
				return err
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
				return err
			}
			zlog.Info("StoreMetadata success", zap.String("scheduler_id", schedulerId), zap.Int64("stock", resp.Affected.Stock), zap.Int64("stock", resp.Affected.Quote))

			return nil
		},
	}
	return p
}
