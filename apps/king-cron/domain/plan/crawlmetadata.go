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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithCrawlMetadata  = "CronWithCrawlMetadata"
	AliasWithCrawlMetadata = "抓取数据"
)

func init() {
	domain.RegisterPlan(NameWithCrawlMetadata, CronWithCrawlMetadata)
}

func CronWithCrawlMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10_SECOND)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithCrawlMetadata, time.Now().Format(time.DateOnly))
			if err == sql.ErrNoRows {
				return domain.Ready, nil
			}
			if err != nil {
				return 0, err
			}

			switch record.Status {
			case domain.ProgressCompleted:
				if record.Code.String == codes.SUCCESS {
					return domain.Completed, nil
				}
				return 0, errors.New(record.ErrorMsg.String)

			case domain.ProgressProcessing:
				return domain.Completed, nil

			default:
				return 0, fmt.Errorf("panic: unknown status, nest status: %v", record.Status)
			}
		},
		Todo: func(schedulerId string) (string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10_SECOND)
			defer cancel()

			_, err := client.DefalutCollector.CrawlMetadataAsync(ctx, &wrapperspb.StringValue{Value: schedulerId})
			return "", err
		},

		WriteToDB: func(schedulerId string, err error) error {
			status, code, errormsg := func() (string, sql.NullString, sql.NullString) {
				if err == nil {
					return domain.ProgressProcessing, sql.NullString{}, sql.NullString{}
				}
				return domain.ProgressCompleted, sql.NullString{String: codes.FAILURE, Valid: true}, sql.NullString{String: err.Error(), Valid: true}
			}()

			now := time.Now()
			record := &db.SchedulerRecord{
				Id:          schedulerId,
				Alias:       AliasWithCrawlMetadata,
				Name:        NameWithCrawlMetadata,
				Date:        now,
				ServiceName: "collector",
				FuncName:    "CrawlMetadataAsync",
				Status:      status,
				Code:        code,
				ErrorMsg:    errormsg,
			}

			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10_SECOND)
			defer cancel()

			if _, err := db.SchedulerRecordWithInsertOne(ctx, mysql.DB, record); err != nil {
				return err
			}
			return nil
		},

		NotifyWithError: func(err error) error {
			return domain.DefaultNotifyWithError(NameWithCrawlMetadata, fmt.Errorf("failure: %v", err), []string{"原始数据", "网络"})
		},

		Status: domain.Ready,
		Name:   NameWithCrawlMetadata,
		Alias:  AliasWithCrawlMetadata,
	}
	return p
}
