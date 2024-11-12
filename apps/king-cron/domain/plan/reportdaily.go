package plan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/eviltomorrow/king/apps/king-cron/domain"
	"github.com/eviltomorrow/king/apps/king-cron/domain/db"
	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
	"github.com/eviltomorrow/king/lib/mathutil"
	"github.com/eviltomorrow/king/lib/setting"
	"github.com/eviltomorrow/king/lib/system"
	"github.com/flosch/pongo2/v6"
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

			status, err := client.DefalutFinder.ReportDaily(context.Background(), &wrapperspb.StringValue{Value: now.Format(time.DateOnly)})
			if err != nil {
				return "", err
			}
			text, err := generateMarketStatusHTMLText(status)
			if err != nil {
				return "", err
			}
			return text, nil
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
			return domain.DefaultNotifyWithMsg(NameWithReportDaily, text, []string{"日报", "今日涨跌幅"})
		},

		Status: domain.Ready,
		Name:   NameWithReportDaily,
		Alias:  AliasWithReportDaily,
	}
}

func generateMarketStatusHTMLText(status *pb.MarketStatus) (string, error) {
	tpl, err := pongo2.FromFile(filepath.Join(system.Directory.EtcDir, "assets", "daily-report.html"))
	if err != nil {
		return "", fmt.Errorf("load daily-report.html failure, nest error: %v", err)
	}

	data := map[string]interface{}{
		"date":              status.Date,
		"shang_zheng_value": status.MarketIndex.ShangZheng.Value,
		"shang_zheng_direction": func() string {
			if status.MarketIndex.ShangZheng.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"shang_zheng_change": status.MarketIndex.ShangZheng.HasChanged,

		"shen_zheng_value": status.MarketIndex.ShenZheng.Value,
		"shen_zheng_direction": func() string {
			if status.MarketIndex.ShenZheng.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"shen_zheng_change": status.MarketIndex.ShenZheng.HasChanged,

		"chuang_ye_value": status.MarketIndex.ChuangYe.Value,
		"chuang_ye_direction": func() string {
			if status.MarketIndex.ChuangYe.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"chuang_ye_change": status.MarketIndex.ChuangYe.HasChanged,

		"ke_chuang50_value": status.MarketIndex.KeChuang_50.Value,
		"ke_chuang50_direction": func() string {
			if status.MarketIndex.KeChuang_50.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"ke_chuang50_change": status.MarketIndex.KeChuang_50.HasChanged,

		"bei_zheng50_value": status.MarketIndex.BeiZheng_50.Value,
		"bei_zheng50_direction": func() string {
			if status.MarketIndex.BeiZheng_50.HasChanged > 0 {
				return "⬆️"
			} else {
				return "⬇️"
			}
		}(),
		"bei_zheng50_change": status.MarketIndex.BeiZheng_50.HasChanged,

		"total":           status.MarketStockCount.Total,
		"rise_gt_7":       status.MarketStockCount.RiseGt_7,
		"rise_gt_7_ratio": mathutil.Trunc4(float64(status.MarketStockCount.RiseGt_7) / float64(status.MarketStockCount.Total) * 100),
		"fell_gt_7":       status.MarketStockCount.RiseGt_7,
		"fell_gt_7_ratio": mathutil.Trunc4(float64(status.MarketStockCount.FellGt_7) / float64(status.MarketStockCount.Total) * 100),
	}

	return tpl.Execute(data)
}
