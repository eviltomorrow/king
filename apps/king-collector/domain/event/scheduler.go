package event

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/conf"
	"github.com/eviltomorrow/king/apps/king-collector/domain/metadata"
	"github.com/eviltomorrow/king/apps/king-collector/domain/notification"
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	lastDay int64 = -1
)

type Scheduler struct {
	config *conf.Collector
}

func NewScheduler(config *conf.Collector) *Scheduler {
	return &Scheduler{config: config}
}

func (s *Scheduler) Run() error {
	c := cron.New()
	_, err := c.AddFunc(s.config.Crontab, func() {
		ctx, span := opentrace.DefaultTracer().Start(context.Background(), "RunCron")
		defer span.End()

		if err := s.fetchMetadataEveryWeekDay(ctx, s.config.Source, s.config.CodeList, s.config.RandomPeriod); err != nil {
			span.SetStatus(codes.Error, "fetchMetadataEveryWeekDay failure")
			span.RecordError(err)
			zlog.Error("Crontab run archive metadata every weekday failure", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}
	c.Start()

	finalizer.RegisterCleanupFuncs(func() error {
		c.Stop()
		return nil
	})
	return nil
}

func (s *Scheduler) fetchMetadataEveryWeekDay(ctx context.Context, source string, baseCodeList []string, randomPeriod []int) error {
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
			if err := notification.SendNtfy(ctx, "Fetch metadata failure", fmt.Sprintf("Possible reason: %v", e)); err != nil {
				span.RecordError(err)
				zlog.Error("Send notification failure", zap.Error(err))
			}

			if err := notification.SendEmail(ctx, "Fetch metadata failure", fmt.Sprintf("Possible reason: %v", e)); err != nil {
				span.RecordError(err)
				zlog.Error("Send email failure", zap.Error(err))
			}
		}
	}()

	ctx, span = opentrace.DefaultTracer().Start(ctx, "SynchronizeDataSlow")
	defer span.End()
	total, ignore, e = metadata.SynchronizeMetadataSlow(ctx, source, baseCodeList, randomPeriod)
	if e != nil {
		span.RecordError(e)
		return e
	}

	ctx, span = opentrace.DefaultTracer().Start(ctx, "ArchiveMetadataToStorage")
	defer span.End()

	_, affectedStock, affectedDay, affectedWeek, e = metadata.ArchiveMetadataToStorage(ctx, begin.Format(time.DateOnly))
	if e != nil {
		span.RecordError(e)
		return e
	}
	if lastDay != -1 {
		if lastDay > total && (lastDay-total) > int64(float64(lastDay)*0.1) {
			e = fmt.Errorf("archive data slow possible missing data, nest last: %v, nest count: %v", lastDay, total)
		}
	}

	lastDay = total
	zlog.Info("Archive metedata every weekday complete", zap.Int64("total", total), zap.Int64("ignore", ignore), zap.Int64("stock-affetced", affectedStock), zap.Int64("day-affected", affectedDay), zap.Int64("week-affected", affectedWeek), zap.Duration("cost", time.Since(begin)))
	return nil
}
