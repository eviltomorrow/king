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

var ArchiveCommand = &cobra.Command{
	Use:   "archive",
	Short: "归档数据",
	Run: func(cmd *cobra.Command, args []string) {
		begin, err := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", beginVar))
		if err != nil {
			log.Printf("[F] 转换开始日期失败, nest error: %v, begin: %s", err, begin)
			return
		}
		end, err := time.Parse(time.DateTime, fmt.Sprintf("%s 23:59:59", endVar))
		if err != nil {
			log.Printf("[F] 转换结束日期失败, nest error: %v, end: %s", err, end)
			return
		}

		stubCollector, closeFunc, err := client.NewCollectorWithTarget(fmt.Sprintf("%s:50003", IPVar))
		if err != nil {
			log.Printf("create collector client failure, nest error: %v", err)
			return
		}
		defer closeFunc()
		ClientCollector = stubCollector

		stubStorage, closeFunc, err := client.NewStorageWithTarget(fmt.Sprintf("%s:50001", IPVar))
		if err != nil {
			log.Printf("create storage client failure, nest error: %v", err)
			return
		}
		defer closeFunc()
		ClientStorage = stubStorage

		for begin.Before(end) {
			var (
				date = begin.Format(time.DateOnly)
				now  = time.Now()
			)

			stocks, days, weeks, err := archiveMetadata(context.Background(), date)
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
	ArchiveCommand.Flags().StringVar(&beginVar, "begin", "", "指定开始日期")
	ArchiveCommand.MarkFlagRequired("begin")

	ArchiveCommand.Flags().StringVar(&endVar, "end", time.Now().Format(time.DateOnly), "指定结束日期")
	ArchiveCommand.Flags().StringVar(&IPVar, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func archiveMetadata(ctx context.Context, date string) (int64, int64, int64, error) {
	target, err := ClientStorage.PushMetadata(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	source, err := ClientCollector.PullMetadata(ctx, &wrapperspb.StringValue{Value: date})
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
