package plan

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/snowflake"
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
				if record.Code == "success" {
					return false, nil
				}
				return false, errors.New(record.ErrorMsg)
			}

			record, err = db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithCrawlMetadata, time.Now().Format(time.DateOnly))
			if err != nil {
				return false, err
			}

			if record.Status == domain.StatusCompleted {
				if record.Code == "success" {
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

			_ = stub

			return nil
		},
	}
	return p
}
