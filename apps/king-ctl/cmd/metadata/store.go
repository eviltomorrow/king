package metadata

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var StoreCommand = &cobra.Command{
	Use:   "store",
	Short: "归档数据",
	Run: func(cmd *cobra.Command, args []string) {
		begin, err := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", begin))
		if err != nil {
			log.Printf("[F] 转换开始日期失败, nest error: %v, begin: %s", err, begin)
			return
		}
		end, err := time.Parse(time.DateTime, fmt.Sprintf("%s 23:59:59", end))
		if err != nil {
			log.Printf("[F] 转换结束日期失败, nest error: %v, end: %s", err, end)
			return
		}

		for begin.Before(end) {
			var (
				date = begin.Format(time.DateOnly)
				now  = time.Now()
			)

			stocks, days, weeks, err := store(context.Background(), date)
			if err != nil {
				log.Printf("归档失败, nest error: %v, date: %v", err, date)
			} else {
				fmt.Printf("归档完成, 日期: %v, 股票数: %v, 日交易数据: %v, 周交易数据: %v, 花费: %v\r\n", date, stocks, days, weeks, time.Since(now))
			}
			begin = begin.Add(24 * time.Hour)
		}
	},
}

func init() {
	StoreCommand.PersistentFlags().StringVar(&begin, "begin", "", "指定开始日期")
	StoreCommand.MarkPersistentFlagRequired("begin")

	StoreCommand.PersistentFlags().StringVar(&end, "end", time.Now().Format(time.DateOnly), "指定结束日期")
	StoreCommand.PersistentFlags().StringVar(&IP, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func store(ctx context.Context, date string) (int64, int64, int64, error) {
	stubStorage, closeFuncStorage, err := client.NewStorageWithTarget(fmt.Sprintf("%s:50001", IP))
	if err != nil {
		return 0, 0, 0, err
	}
	defer closeFuncStorage()

	target, err := stubStorage.PushMetadata(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	stubCollector, closeFuncCollector, err := client.NewCollectorWithTarget(fmt.Sprintf("%s:50003", IP))
	if err != nil {
		return 0, 0, 0, err
	}
	defer closeFuncCollector()

	source, err := stubCollector.FetchMetadata(ctx, &wrapperspb.StringValue{Value: date})
	if err != nil {
		return 0, 0, 0, err
	}
	for {
		md, err := source.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, 0, err
		}

		if err := target.Send(md); err != nil {
			return 0, 0, 0, err
		}
	}

	resp, err := target.CloseAndRecv()
	if err != nil {
		return 0, 0, 0, err
	}
	return resp.Affected.Stocks, resp.Affected.Days, resp.Affected.Weeks, nil
	// return 0, 0, nil
}
