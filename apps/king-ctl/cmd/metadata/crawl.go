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
	Short: "抓取数据",
	Run: func(cmd *cobra.Command, args []string) {
		stub, closeFunc, err := client.NewCollectorWithTarget(fmt.Sprintf("%s:50003", IPVar))
		if err != nil {
			log.Printf("create collector client failure, nest error: %v", err)
			return
		}
		defer closeFunc()
		ClientCollector = stub

		if err := crawlMetadataSync(context.Background()); err != nil {
			log.Printf("crawl data failure, nest error: %v", err)
			return
		}
	},
}

var source string

func init() {
	CrawlCommand.Flags().StringVar(&source, "source", "sina", "指定数据源[sina]")
	CrawlCommand.Flags().StringVar(&IPVar, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func crawlMetadataSync(ctx context.Context) error {
	begin := time.Now()

	resp, err := ClientCollector.CrawlMetadataSync(ctx, &wrapperspb.StringValue{Value: source})
	if err != nil {
		return err
	}
	fmt.Printf("[Status] 完成, 源: %s, 总数: %d, 忽略: %d, 实际数量: %d, 花费: %v\r\n", source, resp.Total, resp.Ignore, resp.Total-resp.Ignore, time.Since(begin))
	return nil
}
