package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/service"
	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/etcd"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/resolver"
)

var (
	beginDateVar string
	endDateVar   string
)

func init() {
	StoreCommand.PersistentFlags().StringVar(&beginDateVar, "begin", "", "specify the begin date")
	StoreCommand.MarkPersistentFlagRequired("begin")

	StoreCommand.PersistentFlags().StringVar(&endDateVar, "end", time.Now().Format(time.DateOnly), "specify the end date")

	RootCommand.AddCommand(StoreCommand)
}

var StoreCommand = &cobra.Command{
	Use:   "store",
	Short: "Store metadata to storage",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, span := opentrace.DefaultTracer().Start(context.Background(), "StoreMetadata")
		defer span.End()

		if err := loadConfig(); err != nil {
			span.SetStatus(codes.Error, "Load config failure")
			span.RecordError(err)
			log.Printf("[F] Load config failure, nest error: %v", err)
			return
		}

		mongodb.DSN = cfg.MongoDB.DSN
		if err := mongodb.Connect(); err != nil {
			span.SetStatus(codes.Error, "New mongodb client failure")
			span.RecordError(err)
			log.Printf("[F] Create mongodb client failure, nest error: %v", err)
			return
		}

		etcd.Endpoints = cfg.Etcd.Endpoints
		client, err := etcd.NewClient()
		if err != nil {
			span.SetStatus(codes.Error, "New etcd client failure")
			span.RecordError(err)
			log.Printf("[F] Create etcd client failure, nest error: %v", err)
			return
		}
		defer client.Close()
		resolver.Register(lb.NewBuilder(client))

		begin, err := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", beginDateVar))
		if err != nil {
			span.SetStatus(codes.Error, "Parse begin failure")
			span.RecordError(err)
			log.Printf("[F] Parse begin date failure, nest error: %v, begin: %s", err, beginDateVar)
			return
		}
		end, err := time.Parse(time.DateTime, fmt.Sprintf("%s 23:59:59", endDateVar))
		if err != nil {
			span.SetStatus(codes.Error, "Parse end failure")
			span.RecordError(err)
			log.Printf("[F] Parse end date failure, nest error: %v, end: %s", err, endDateVar)
			return
		}

		for begin.Before(end) {
			var (
				date = begin.Format(time.DateOnly)
				now  = time.Now()
			)
			stock, day, week, err := service.StoreMetadataToStorage(ctx, date)
			if err != nil {
				log.Printf("[E] Store data failure, nest error: %v, date: %v", err, date)
			} else {
				fmt.Printf("[I] Store data success, date: %v, affected-stock: %v, affected-day: %v, affetced-week: %v, cost: %v\r\n", date, stock, day, week, time.Since(now))
			}
			begin = begin.Add(24 * time.Hour)
		}
	},
}
