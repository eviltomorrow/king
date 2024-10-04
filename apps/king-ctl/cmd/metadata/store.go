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

			stock, quote, err := store(context.Background(), date)
			if err != nil {
				log.Printf("归档失败, nest error: %v, date: %v", err, date)
			} else {
				fmt.Printf("归档完成, 日期: %v, 股票数量: %v, 股票交易数量: %v, 花费: %v\r\n", date, stock, quote, time.Since(now))
			}
			begin = begin.Add(24 * time.Hour)
		}
	},
}

var (
	begin string
	end   string
)

func init() {
	StoreCommand.PersistentFlags().StringVar(&begin, "begin", "", "指定开始日期")
	StoreCommand.MarkPersistentFlagRequired("begin")

	StoreCommand.PersistentFlags().StringVar(&end, "end", time.Now().Format(time.DateOnly), "指定结束日期")
	StoreCommand.PersistentFlags().StringVar(&IP, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func store(ctx context.Context, date string) (int64, int64, error) {
	stub, closeFunc, err := client.NewCollectorWithTarget(fmt.Sprintf("%s:50003", IP))
	if err != nil {
		return 0, 0, err
	}
	defer closeFunc()

	resp, err := stub.StoreMetadata(ctx, &wrapperspb.StringValue{Value: date})
	if err != nil {
		return 0, 0, err
	}
	return resp.Affected.Stock, resp.Affected.Quote, nil
}
