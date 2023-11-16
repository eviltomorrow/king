package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/db/mongodb"
	"github.com/eviltomorrow/king/lib/etcd"
	grpcclient "github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
		ctx, span := opentrace.DefaultTracer().Start(context.Background(), "Manual StoreMetadata")
		defer span.End()

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

		for begin.Before(end) {
			var (
				date = begin.Format(time.DateOnly)
				now  = time.Now()
			)
			total, stock, day, week, err := store(ctx, date)
			if err != nil {
				log.Printf("[E] => Store data failure, nest error: %v, date: %v", err, date)
			} else {
				fmt.Printf("[I] => Store data success, date: %v, total: %v, stock: %v, day: %v, week: %v, cost: %v\r\n", date, total, stock, day, week, time.Since(now))
			}
			begin = begin.Add(24 * time.Hour)
		}
	},
}

func store(ctx context.Context, date string) (int64, int64, int64, int64, error) {
	stub, closeFunc, err := grpcclient.NewCollectorWithEtcd()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer closeFunc()

	resp, err := stub.StoreMetadata(ctx, &wrapperspb.StringValue{Value: date})
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return resp.Total, resp.Stock, resp.Day, resp.Week, nil
}
