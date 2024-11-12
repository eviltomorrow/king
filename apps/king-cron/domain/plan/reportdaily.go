package plan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/apps/king-cron/domain/notification"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/grpc/transformer"
	"github.com/eviltomorrow/king/lib/setting"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NameWithReportDaily  = "CronWithReportDaily"
	AliasWithReportDaily = "日报"
)

func init() {
	domain.RegisterPlan(AliasWithReportDaily, CronWithReportDaily)
}

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

			status, err := client.DefaultFinder.ReportDaily(context.Background(), &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
			if err != nil {
				return "", err
			}
			data := transformer.GenerateMarketStatusToMap(status)

			value := make(map[string]string)
			for k, v := range data {
				value[k] = fmt.Sprintf("%v", v)
			}
			buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(value)
			if err != nil {
				return "", err
			}
			return string(buf), err
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
				Alias:       AliasWithReportDaily,
				Name:        NameWithReportDaily,
				Date:        now,
				ServiceName: "brain",
				FuncName:    "ReportDaily",
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

		NotifyWithData: func(text string) error {
			data := make(map[string]string)
			if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(text), &data); err != nil {
				return err
			}

			resp, err := client.DefaultTemplate.Render(context.Background(), &pb.RenderRequest{
				TemplateName: "daily_report.html",
				Data:         data,
			})
			if err != nil {
				return err
			}
			if err := notification.DefaultNotifyForEmailWithMsg(NameWithReportDaily, resp.Value); err != nil {
				return err
			}

			resp, err = client.DefaultTemplate.Render(context.Background(), &pb.RenderRequest{
				TemplateName: "daily_report.txt",
				Data:         data,
			})
			if err != nil {
				return err
			}
			if err := notification.DefaultNotifyForNtfyWithMsg(fmt.Sprintf("%s 日 汇总", time.Now().Format(time.DateOnly)), resp.Value, []string{"简报", "股票", "统计"}); err != nil {
				return err
			}
			return nil
		},

		Status: domain.Ready,
		Name:   NameWithReportDaily,
		Alias:  AliasWithReportDaily,
	}
}
