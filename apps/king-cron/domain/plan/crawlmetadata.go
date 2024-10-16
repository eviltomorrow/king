package plan

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/setting"
	"github.com/eviltomorrow/king/lib/snowflake"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithCrawlMetadata = "爬取数据"
)

func CronWithCrawlMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: nil,
		Todo: func() (string, error) {
			stub, shutdown, err := client.NewCollectorWithEtcd()
			if err != nil {
				return "", err
			}
			defer shutdown()

			schedulerId := snowflake.GenerateID()

			ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10SECOND)
			defer cancel()

			md := metadata.MD{
				"scheduler_id": []string{schedulerId},
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
			if _, err := stub.CrawlMetadataAsync(ctx, &wrapperspb.StringValue{Value: "sina"}); err != nil {
				return "", err
			}

			t := time.Now().Format(time.DateOnly)
			now, err := time.Parse(time.DateOnly, t)
			if err != nil {
				return "", err
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
				return "", err
			}

			return "", nil
		},

		NotifyWithError: func(err error) error {
			return domain.DefaultNotifyWithError(NameWithCrawlMetadata, fmt.Errorf("failure: %v", err), []string{"原始数据", "网络"})
		},

		Status: domain.Ready,
		Name:   NameWithCrawlMetadata,
	}
	return p
}
