package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var CrawlCommand = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl data manual from specify source[sina/net126]",
	Run: func(cmd *cobra.Command, args []string) {
		if err := crawl(context.Background()); err != nil {
			log.Printf("crawl data failure, nest error: %v", err)
			return
		}
	},
}

var (
	source string
)

func init() {
	CrawlCommand.PersistentFlags().StringVar(&source, "source", "sina", "crawl data from [sina/net126]")
}

func crawl(ctx context.Context) error {
	begin := time.Now()

	stub, closeFunc, err := client.NewCollectorWithEtcd()
	if err != nil {
		return err
	}
	defer closeFunc()

	resp, err := stub.CrawlMetadata(ctx, &wrapperspb.StringValue{Value: source})
	if err != nil {
		return err
	}
	fmt.Printf("[Status] Complete, Source: %s, Total: %d, Ignore: %d, Actual: %d, Cost: %v\r\n", source, resp.Total, resp.Ignore, resp.Total-resp.Ignore, time.Since(begin))
	return nil
}
