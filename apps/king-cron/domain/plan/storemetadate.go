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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithStoreMetadata  = "CronWithStoreMetadata"
	AliasWithStoreMetadata = "存储数据"
)

func init() {
	domain.RegisterPlan(NameWithStoreMetadata, CronWithStoreMetadata)
}

func CronWithStoreMetadata() *domain.Plan {
	p := &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10TIMEOUT)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithStoreMetadata, time.Now().Format(time.DateOnly))
			if err != nil && err != sql.ErrNoRows {
				return 0, err
			}

			if record != nil && record.Status == domain.ProgressCompleted {
				if record.Code.String == codes.SUCCESS {
					return domain.Completed, nil
				}
				return 0, errors.New(record.ErrorMsg.String)
			}

			record, err = db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithCrawlMetadata, time.Now().Format(time.DateOnly))
			if err != nil && err != sql.ErrNoRows {
				return 0, err
			}

			if record != nil && record.Status == domain.ProgressCompleted {
				if record.Code.String == codes.SUCCESS {
					return domain.Ready, nil
				}
				return 0, errors.New(record.ErrorMsg.String)
			}
			return domain.Pending, nil
		},

		Todo: func(schedulerId string) (string, error) {
			target, err := client.DefalutStorage.PushMetadata(context.Background())
			if err != nil {
				return "", err
			}

			now := time.Now()

			source, err := client.DefalutCollector.FetchMetadata(context.Background(), &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
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

			if _, err := target.CloseAndRecv(); err != nil {
				return "", err
			}

			return "", nil
		},

		WriteToDB: func(schedulerId string, err error) error {
			status, code, errormsg := func() (string, sql.NullString, sql.NullString) {
				if err == nil {
					return domain.ProgressCompleted, sql.NullString{String: codes.SUCCESS, Valid: true}, sql.NullString{}
				}
				return domain.ProgressCompleted, sql.NullString{String: codes.FAILURE, Valid: true}, sql.NullString{String: err.Error(), Valid: true}
			}()

			now := time.Now()
			record := &db.SchedulerRecord{
				Id:          schedulerId,
				Alias:       AliasWithStoreMetadata,
				Name:        NameWithStoreMetadata,
				Date:        now,
				ServiceName: "storage",
				FuncName:    "PushMetadata",
				Status:      status,
				Code:        code,
				ErrorMsg:    errormsg,
			}

			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10TIMEOUT)
			defer cancel()

			if _, err := db.SchedulerRecordWithInsertOne(ctx, mysql.DB, record); err != nil {
				return err
			}
			return nil
		},

		NotifyWithError: func(err error) error {
			return domain.DefaultNotifyWithError(NameWithStoreMetadata, fmt.Errorf("failure: %v", err), []string{"缓存数据", "数据库"})
		},

		Status: domain.Ready,
		Name:   NameWithStoreMetadata,
		Alias:  AliasWithStoreMetadata,
	}
	return p
}
