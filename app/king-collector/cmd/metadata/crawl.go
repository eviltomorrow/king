package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/etcd"
	grpcclient "github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/grpc/lb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	sourceVar string
)

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
			log.Fatalf("[F] Load config failure, nest error: %v", err)
		}

		etcd.Endpoints = cfg.Etcd.Endpoints
		client, err := etcd.NewClient()
		if err != nil {
			log.Fatalf("[F] Create etcd client failure, nest error: %v", err)
		}
		defer client.Close()

		resolver.Register(lb.NewBuilder(client))

		if err := crawl(); err != nil {
			log.Fatalf("[F] Crawl data failure, nest error: %v", err)
		}
	},
}

func crawl() error {
	var begin = time.Now()

	stub, closeFunc, err := grpcclient.NewCollectorWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	resp, err := stub.CrawlMetadata(context.Background(), &wrapperspb.StringValue{Value: sourceVar})
	if err != nil {
		return err
	}
	fmt.Printf("[Status] Complete, Source: %s, Total: %d, Ignore: %d, Actual: %d, Cost: %v\r\n", sourceVar, resp.Total, resp.Ignore, resp.Total-resp.Ignore, time.Since(begin))
	return nil
}
