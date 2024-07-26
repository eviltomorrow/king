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

var ArchiveCommand = &cobra.Command{
	Use:   "archive",
	Short: "Archive metadata to storage",
	Run: func(cmd *cobra.Command, args []string) {
		begin, err := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", begin))
		if err != nil {
			log.Printf("[F] Parse begin date failure, nest error: %v, begin: %s", err, begin)
			return
		}
		end, err := time.Parse(time.DateTime, fmt.Sprintf("%s 23:59:59", end))
		if err != nil {
			log.Printf("[F] Parse end date failure, nest error: %v, end: %s", err, end)
			return
		}

		for begin.Before(end) {
			var (
				date = begin.Format(time.DateOnly)
				now  = time.Now()
			)

			total, stock, day, week, err := store(context.Background(), date)
			if err != nil {
				log.Printf("Archive metadata failure, nest error: %v, date: %v", err, date)
			} else {
				fmt.Printf("Archive metadata success, date: %v, total: %v, stock: %v, day: %v, week: %v, cost: %v\r\n", date, total, stock, day, week, time.Since(now))
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
	ArchiveCommand.PersistentFlags().StringVar(&begin, "begin", "", "specify the begin date")
	ArchiveCommand.MarkPersistentFlagRequired("begin")

	ArchiveCommand.PersistentFlags().StringVar(&end, "end", time.Now().Format(time.DateOnly), "specify the end date")

}

func store(ctx context.Context, date string) (int64, int64, int64, int64, error) {
	stub, closeFunc, err := client.NewCollectorWithEtcd()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer closeFunc()

	resp, err := stub.ArchiveMetadata(ctx, &wrapperspb.StringValue{Value: date})
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return resp.Total, resp.Affected.Stock, resp.Affected.QuoteDay, resp.Affected.QuoteWeek, nil
}
