package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/service/db"
	"github.com/eviltomorrow/king/apps/king-collector/service/synchronize"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	grpcclient "github.com/eviltomorrow/king/lib/grpc/client"
	emailpb "github.com/eviltomorrow/king/lib/grpc/pb/king-email"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	Source        = "sina"
	lastDay int64 = -1
)

func FetchMetadataEveryWeekDay(ctx context.Context) error {
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

			ctx, span = opentrace.DefaultTracer().Start(ctx, "NotifyWithEmail")
			defer span.End()
			if err := NotifyWithEmail(ctx, fmt.Sprintf("Possible reason: %v", e)); err != nil {
				span.RecordError(err)
				zlog.Error("Notify with email failure", zap.Error(err))
			}
		}
	}()

	ctx, span = opentrace.DefaultTracer().Start(ctx, "DataSlow")
	defer span.End()
	total, ignore, e = synchronize.DataSlow(Source)
	if e != nil {
		span.RecordError(e)
		return e
	}

	ctx, span = opentrace.DefaultTracer().Start(ctx, "StoreMetadataToStorage")
	defer span.End()
	_, affectedStock, affectedDay, affectedWeek, e = StoreMetadataToStorage(ctx, begin.Format(time.DateOnly))
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

func StoreMetadataToStorage(ctx context.Context, date string) (int64, int64, int64, int64, error) {
	var span trace.Span
	ctx, span = opentrace.DefaultTracer().Start(ctx, "StoreMetadataToStorage")
	defer span.End()

	client, closeFunc, err := grpcclient.NewStorageWithEtcd()
	if err != nil {
		span.RecordError(err)
		return 0, 0, 0, 0, err
	}
	defer closeFunc()

	stub, err := client.PushMetadata(ctx)
	if err != nil {
		span.RecordError(err)
		return 0, 0, 0, 0, err
	}

	var (
		offset, limit, total int64 = 0, 100, 0
		lastID               string
		timeout              = 20 * time.Second
	)

	for {
		metadata, err := db.SelectMetadataRange(mongodb.DB, offset, limit, date, lastID, timeout)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		for _, md := range metadata {
			if err := stub.Send(&pb.Metadata{
				Source:          md.Source,
				Code:            md.Code,
				Name:            md.Name,
				Open:            md.Open,
				YesterdayClosed: md.YesterdayClosed,
				Latest:          md.Latest,
				High:            md.High,
				Low:             md.Low,
				Volume:          md.Volume,
				Account:         md.Account,
				Date:            md.Date,
				Time:            md.Time,
				Suspend:         md.Suspend,
			}); err != nil {
				return 0, 0, 0, 0, err
			}
			total++
		}
		if len(metadata) < int(limit) {
			break
		}
		offset += limit
	}

	resp, err := stub.CloseAndRecv()
	if err != nil {
		span.SetStatus(codes.Error, "CloseAndRecv result failure")
		span.RecordError(err)
		return 0, 0, 0, 0, err
	}
	return total, resp.StockAffected, resp.QuoteDayAffected, resp.QuoteWeekAffected, nil
}

func NotifyWithEmail(ctx context.Context, reason string) error {
	client, closeFunc, err := grpcclient.NewEmailWithEtcd()
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
			"content1": "King-collector store data failure!",
			"content2": reason,
			"end":      "Bye",
		},
	})
	if err != nil {
		return err
	}

	if _, err = client.Send(ctx, &emailpb.Mail{
		To: []*emailpb.Contact{
			{Name: "Shepard", Address: "eviltomorrow@163.com"},
		},
		Subject: fmt.Sprintf("(%s): King-collector Store Data Failure", time.Now().Format(time.DateOnly)),
		Body:    msg.Value,
	}); err != nil {
		return err
	}

	return nil
}
