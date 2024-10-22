package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	"github.com/eviltomorrow/king/lib/timeutil"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var StatsCommand = &cobra.Command{
	Use:   "stats",
	Short: "统计数据",
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
			)

			status := BoldGreen.Sprint("正常")
			days, weeks, err := stats(context.Background(), date)
			if err != nil {
				log.Printf("数据统计失败, nest error: %v, date: %v", err, date)
			} else {
				if days == 0 && (begin.Weekday() == time.Monday ||
					begin.Weekday() == time.Tuesday ||
					begin.Weekday() == time.Wednesday ||
					begin.Weekday() == time.Thursday ||
					begin.Weekday() == time.Friday) {
					status = BoldRed.Sprint("缺失")
				}
				fmt.Printf("=> 数据统计, 日期: %v(%v), 日交易数据: %4v, 周交易数据: %4v, 状态: (%v)\r\n", date, timeutil.ParseWeekday(begin), days, weeks, status)
			}
			begin = begin.Add(24 * time.Hour)
		}
	},
}

var (
	BoldGreen = color.New(color.FgGreen, color.Bold)
	BoldRed   = color.New(color.FgRed, color.Bold)
)

func init() {
	StatsCommand.PersistentFlags().StringVar(&begin, "begin", "", "指定开始日期")
	StatsCommand.MarkPersistentFlagRequired("begin")

	StatsCommand.PersistentFlags().StringVar(&end, "end", time.Now().Format(time.DateOnly), "指定结束日期")
	StatsCommand.PersistentFlags().StringVar(&IP, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func stats(ctx context.Context, date string) (int64, int64, error) {
	stub, closeFunc, err := client.NewStorageWithTarget(fmt.Sprintf("%s:50001", IP))
	if err != nil {
		return 0, 0, err
	}
	defer closeFunc()

	resp, err := stub.StatsMetadata(ctx, &wrapperspb.StringValue{Value: date})
	if err != nil {
		return 0, 0, err
	}

	return resp.Queried.Days, resp.Queried.Weeks, nil
	// return 0, 0, nil
}
