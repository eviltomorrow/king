package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
	"github.com/eviltomorrow/king/lib/timeutil"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ShowCommand = &cobra.Command{
	Use:   "show",
	Short: "统计数据",
	Run: func(cmd *cobra.Command, args []string) {
		end, err := time.Parse(time.DateTime, fmt.Sprintf("%s 23:59:59", endVar))
		if err != nil {
			log.Printf("[F] 转换结束日期失败, nest error: %v, end: %s", err, end)
			return
		}

		begin := end.Add(-30 * 25 * time.Hour)
		if beginVar != "" {
			begin, err = time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", beginVar))
			if err != nil {
				log.Printf("[F] 转换开始日期失败, nest error: %v, begin: %s", err, begin)
				return
			}
		}

		stub, closeFunc, err := client.NewStorageWithTarget(fmt.Sprintf("%s:50001", IPVar))
		if err != nil {
			log.Printf("create storage client failure, nest error: %v", err)
			return
		}
		defer closeFunc()
		ClientStorage = stub

		for begin.Before(end) {
			date := begin.Format(time.DateOnly)
			var (
				days, weeks int64
				err         error
			)
			status := BoldGreen.Sprint("正常")

			days, weeks, err = showQuote(context.Background(), begin)
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
	ShowCommand.Flags().StringVar(&beginVar, "begin", "", "指定开始日期")
	ShowCommand.Flags().StringVar(&endVar, "end", time.Now().Format(time.DateOnly), "指定结束日期")
	ShowCommand.Flags().StringVar(&IPVar, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func showQuote(ctx context.Context, date time.Time) (int64, int64, error) {
	day, err := countQuoteDay(ctx, date.Format(time.DateOnly))
	if err != nil {
		return 0, 0, err
	}

	var week int64
	if date.Weekday() == time.Friday {
		week, err = countQuoteWeek(ctx, date.Format(time.DateOnly))
		if err != nil {
			return 0, 0, err
		}
	}
	return day, week, nil
}

func countQuoteDay(ctx context.Context, date string) (int64, error) {
	resp, err := ClientStorage.CountQuote(ctx, &pb.CountQuoteRequest{Date: date, Kind: pb.CountQuoteRequest_Day})
	if err != nil {
		return 0, err
	}
	return resp.Value, nil
}

func countQuoteWeek(ctx context.Context, date string) (int64, error) {
	resp, err := ClientStorage.CountQuote(ctx, &pb.CountQuoteRequest{Date: date, Kind: pb.CountQuoteRequest_Week})
	if err != nil {
		return 0, err
	}
	return resp.Value, nil
}
