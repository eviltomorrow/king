package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/app/king-collector/service/synchronize"
	grpcclient "github.com/eviltomorrow/king/lib/grpc/client"
	emailpb "github.com/eviltomorrow/king/lib/grpc/pb/king-email"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	Source        = "sina"
	lastDay int64 = -1
)

func ArchiveMetadataEveryWeekDay() error {
	zlog.Info("Archive metadata every weekday begin")
	var (
		affectedStock, affectedDay, affectedWeek int64
		total, ignore                            int64
		e                                        error
		begin                                    = time.Now()
	)

	defer func() {
		if e != nil {
			zlog.Error("Archive metadata every weekday has something wrong", zap.Error(e))
			if err := notifyWithEmail(fmt.Sprintf("Possible reason: %v", e)); err != nil {
				zlog.Error("Notify with email failure", zap.Error(err))
			}
		}
	}()

	total, ignore, e = synchronize.DataSlow(Source)
	if e != nil {
		return e
	}

	affectedStock, affectedDay, affectedWeek, e = ArchiveMetadataToRepository(begin.Format(time.DateOnly))
	if e != nil {
		return e
	}
	if lastDay != -1 {
		if lastDay > total && (lastDay-total) > int64(float64(lastDay)*0.1) {
			e = fmt.Errorf("synchronize data slow possible missing data, nest last: %v, nest count: %v", lastDay, total)
		}
	}

	lastDay = total
	zlog.Info("Arichive metedata every weekday complete", zap.Int64("synchronize-total", total), zap.Int64("synchronize-ignore", ignore), zap.Int64("affetced-stock", affectedStock), zap.Int64("affected-day", affectedDay), zap.Int64("affected-week", affectedWeek), zap.Duration("cost", time.Since(begin)))
	return nil
}

func ArchiveMetadataToRepository(date string) (int64, int64, int64, error) {
	client, closeFunc, err := grpcclient.NewRepository()
	if err != nil {
		return 0, 0, 0, err
	}
	defer closeFunc()

	resp, err := client.ArchiveMetadata(context.Background(), &wrapperspb.StringValue{Value: date})
	if err != nil {
		return 0, 0, 0, err
	}

	return resp.AffectedStock, resp.AffectedQuoteDay, resp.AffectedQuoteWeek, nil
}

func notifyWithEmail(reason string) error {
	client, closeFunc, err := grpcclient.NewEmail()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg, err := client.RenderTemplate(ctx, &emailpb.TemplateData{
		Name: "",
		Data: map[string]string{
			"user":     "Shepard",
			"content1": "King-collector synchronize data failure!",
			"content2": reason,
			"end":      "Bye",
		},
	})
	if err != nil {
		return err
	}

	_, err = client.Send(ctx, &emailpb.Mail{
		To: []*emailpb.Contact{
			{Name: "Shepard", Address: "eviltomorrow@163.com"},
		},
		Subject: fmt.Sprintf("(%s): King-collector Synchronize Data Failure", time.Now().Format(time.DateOnly)),
		Body:    msg.Value,
	})
	return err
}
