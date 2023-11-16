package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/etcd"
	grpcclient "github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/eviltomorrow/king/lib/opentrace"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var sourceVar string

func init() {
	CrawlCommand.PersistentFlags().StringVar(&sourceVar, "source", "sina", "crawl data from [sina/net126]")
	CrawlCommand.MarkPersistentFlagRequired("source")

	RootCommand.AddCommand(CrawlCommand)
}

var CrawlCommand = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl data manual from specify source[sina/net126]",
	Run: func(cmd *cobra.Command, args []string) {
		if err := loadConfig(); err != nil {
			log.Printf("[F] Load config failure, nest error: %v", err)
			return
		}
		opentrace.OtelDSN = cfg.Otel.DSN
		destroy, err := opentrace.InitTraceProvider()
		if err != nil {
			log.Printf("[F] Init trace provider failure, nest error: %v", err)
			return
		}
		defer destroy()

		ctx, span := opentrace.DefaultTracer().Start(context.Background(), "Manual crawl metadata")
		defer span.End()

		etcd.Endpoints = cfg.Etcd.Endpoints
		client, err := etcd.NewClient()
		if err != nil {
			log.Printf("[F] Create etcd client failure, nest error: %v", err)
			span.RecordError(err)
			return
		}
		defer client.Close()

		resolver.Register(lb.NewBuilder(client))

		if err := crawl(ctx); err != nil {
			log.Printf("[F] Crawl data failure, nest error: %v", err)
			span.RecordError(err)
			return
		}
	},
}

func crawl(ctx context.Context) error {
	begin := time.Now()

	stub, closeFunc, err := grpcclient.NewCollectorWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	ctx, span := opentrace.DefaultTracer().Start(ctx, "CrawlMetadata")
	defer span.End()

	span.SetAttributes(attribute.String("source", sourceVar))

	resp, err := stub.CrawlMetadata(ctx, &wrapperspb.StringValue{Value: sourceVar})
	if err != nil {
		span.RecordError(err)
		return err
	}
	fmt.Printf("[Status] Complete, Source: %s, Total: %d, Ignore: %d, Actual: %d, Cost: %v\r\n", sourceVar, resp.Total, resp.Ignore, resp.Total-resp.Ignore, time.Since(begin))
	return nil
}
