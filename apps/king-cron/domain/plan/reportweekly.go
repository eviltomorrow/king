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
	"github.com/eviltomorrow/king/lib/grpc/transformer"
	"github.com/eviltomorrow/king/lib/setting"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithReportWeekly  = "CronWithReportWeekly"
	AliasWithReportWeekly = "周报"
)

func init() {
	domain.RegisterPlan(NameWithReportWeekly, CronWithReportWeekly)
}

func CronWithReportWeekly() *domain.Plan {
	return &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10_SECOND)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithReportWeekly, time.Now().Format(time.DateOnly))
			if err != nil && err != sql.ErrNoRows {
				return 0, err
			}

			if record != nil && record.Status == domain.ProgressCompleted {
				if record.Code.String == codes.SUCCESS {
					return domain.Completed, nil
				}
				return 0, errors.New(record.ErrorMsg.String)
			}

			record, err = db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithStoreMetadata, time.Now().Format(time.DateOnly))
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
		Todo: func(string) error {
			now := time.Now()

			status, err := client.DefaultFinder.ReportWeekly(context.Background(), &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
			if err != nil {
				return err
			}
			data := transformer.GenerateMarketStatusToMap(status)

			value := make(map[string]string)
			for k, v := range data {
				value[k] = fmt.Sprintf("%v", v)
			}

			e := make([]error, 0, 2)
			if err := domain.NotifyForEmail("weekly-report.html", fmt.Sprintf("%s 日 汇总", time.Now().Format(time.DateOnly)), value); err != nil {
				zlog.Error("Notify for email failure", zap.Error(err))
				e = append(e, err)
			}

			if err := domain.NotifyForNtfy("weekly-report.txt", fmt.Sprintf("%s 日 汇总", time.Now().Format(time.DateOnly)), value); err != nil {
				zlog.Error("Notify for ntfy failure", zap.Error(err))
				e = append(e, err)
			}

			if len(e) == 2 {
				return errors.Join(e...)
			}
			return nil
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
				Alias:       AliasWithReportWeekly,
				Name:        NameWithReportWeekly,
				Date:        now,
				ServiceName: "brain",
				FuncName:    "ReportWeekly",
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

		Status: domain.Ready,
		Name:   NameWithReportWeekly,
		Alias:  AliasWithReportWeekly,
	}
}
