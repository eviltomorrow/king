package event

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/domain/service"
	client_grpc "github.com/eviltomorrow/king/lib/grpc/client"
	pb_notification "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	Source        = "sina"
	lastDay int64 = -1
)

func TriggerFetchMetadataEveryWeekDay(ctx context.Context) error {
	var span trace.Span
	ctx, span = opentrace.DefaultTracer().Start(ctx, "FetchMetadataEveryWeekDay")
	defer span.End()

	zlog.Info("Fetch metadata every weekday begin")
	var (
		affectedStock, affectedDay, affectedWeek int64
		total, ignore                            int64
		e                                        error
		begin                                    = time.Now()
	)

	defer func() {
		if e != nil {
			zlog.Error("Fetch metadata every weekday has something wrong", zap.Error(e))

			ctx, span = opentrace.DefaultTracer().Start(ctx, "Send notification")
			defer span.End()
			if err := SendNotification(ctx, "Fetch metadata failure", fmt.Sprintf("Possible reason: %v", e)); err != nil {
				span.RecordError(err)
				zlog.Error("Send notification failure", zap.Error(err))
			}

			if err := SendEmail(ctx, "Fetch metadata failure", fmt.Sprintf("Possible reason: %v", e)); err != nil {
				span.RecordError(err)
				zlog.Error("Send email failure", zap.Error(err))
			}
		}
	}()

	ctx, span = opentrace.DefaultTracer().Start(ctx, "SynchronizeDataSlow")
	defer span.End()
	total, ignore, e = service.SynchronizeMetadataSlow(Source)
	if e != nil {
		span.RecordError(e)
		return e
	}

	ctx, span = opentrace.DefaultTracer().Start(ctx, "StoreMetadataToStorage")
	defer span.End()

	_, affectedStock, affectedDay, affectedWeek, e = service.PushMetadataToStorage(ctx, begin.Format(time.DateOnly))
	if e != nil {
		span.RecordError(e)
		return e
	}
	if lastDay != -1 {
		if lastDay > total && (lastDay-total) > int64(float64(lastDay)*0.1) {
			e = fmt.Errorf("store data slow possible missing data, nest last: %v, nest count: %v", lastDay, total)
		}
	}

	lastDay = total
	zlog.Info("Store metedata every weekday complete", zap.Int64("total", total), zap.Int64("ignore", ignore), zap.Int64("stock-affetced", affectedStock), zap.Int64("day-affected", affectedDay), zap.Int64("week-affected", affectedWeek), zap.Duration("cost", time.Since(begin)))
	return nil
}

func SendEmail(ctx context.Context, subject, reason string) error {
	client, closeFunc, err := client_grpc.NewEmailWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err = client.Send(ctx, &pb_notification.Mail{
		To: []*pb_notification.Contact{
			{Name: "Shepard", Address: "eviltomorrow@163.com"},
		},
		Subject: subject,
		Body:    reason,
	}); err != nil {
		return err
	}

	return nil
}

func SendNotification(ctx context.Context, title, reason string) error {
	client, closeFunc, err := client_grpc.NewNotificationWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err = client.Send(ctx, &pb_notification.Msg{
		Topic:    "topic_alert",
		Message:  reason,
		Title:    title,
		Priority: 4,
		Tags:     []string{"warning", "metadata", "crawl"},
	}); err != nil {
		return err
	}

	return nil
}
