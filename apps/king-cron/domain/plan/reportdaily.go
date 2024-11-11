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
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"github.com/eviltomorrow/king/lib/setting"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithReportDaily  = "CronWithReportDaily"
	AliasWithReportDaily = "日报"
)

func CronWithReportDaily() *domain.Plan {
	return &domain.Plan{
		Precondition: func() (domain.StatusCode, error) {
			ctx, cancel := context.WithTimeout(context.Background(), setting.DEFUALT_HANDLE_10_SECOND)
			defer cancel()

			record, err := db.SchedulerRecordWithSelectOneByDateName(ctx, mysql.DB, NameWithReportDaily, time.Now().Format(time.DateOnly))
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
		Todo: func(string) (string, error) {
			now := time.Now()

			status, err := client.DefalutFinder.ReportDaily(context.Background(), &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
			if err != nil {
				return "", err
			}
			text, err := generateHTMLText(status)
			if err != nil {
				return "", err
			}
			return text, nil
		},
		WriteToDB: func(string, error) error {
			return nil
		},

		NotifyWithData: func(string) error {
			return nil
		},

		Status: domain.Ready,
		Name:   NameWithReportDaily,
		Alias:  AliasWithReportDaily,
	}
}

func generateHTMLText(status *pb.MarketStatus) (string, error) {
	_ = status
	return "", nil
}
