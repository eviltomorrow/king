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
	"github.com/eviltomorrow/king/lib/setting"
	"github.com/eviltomorrow/king/lib/snowflake"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithCrawlMetadata = "爬取数据"
)

func CronWithCrawlMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10TIMEOUT)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithCrawlMetadata, time.Now().Format(time.DateOnly))
			if err != nil && err != sql.ErrNoRows {
				return 0, err
			}
			if record != nil && record.Status == domain.StatusCompleted {
				if record.Code.String == codes.SUCCESS {
					return domain.Completed, nil
				}
				return 0, errors.New(record.ErrorMsg.String)
			}

			return domain.Ready, nil
		},
		Todo: func() (string, error) {
			stub, shutdown, err := client.NewCollectorWithEtcd()
			if err != nil {
				return "", err
			}
			defer shutdown()

			schedulerId := snowflake.GenerateID()
			ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10SECOND)
			defer cancel()

			_, err = stub.CrawlMetadataAsync(ctx, &wrapperspb.StringValue{Value: schedulerId})
			return "", err
		},

		NotifyWithError: func(err error) error {
			return domain.DefaultNotifyWithError(NameWithCrawlMetadata, fmt.Errorf("failure: %v", err), []string{"原始数据", "网络"})
		},

		CallFuncInfo: domain.CallFuncInfo{
			ServiceName: "collector",
			FuncName:    "CrawlMetadataAsync",
		},
		Type:   domain.ASYNC,
		Status: domain.Ready,
		Name:   NameWithCrawlMetadata,
	}
	return p
}
