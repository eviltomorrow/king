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
	beginDateVar string
	endDateVar   string
)

func init() {
	ArchiveCommand.PersistentFlags().StringVar(&beginDateVar, "begin", "", "specify the begin date")
	ArchiveCommand.MarkPersistentFlagRequired("begin")

	ArchiveCommand.PersistentFlags().StringVar(&endDateVar, "end", "", "specify the end date")
	ArchiveCommand.MarkPersistentFlagRequired("end")

	RootCommand.AddCommand(ArchiveCommand)
}

var ArchiveCommand = &cobra.Command{
	Use:   "archive",
	Short: "Archive metadata to storage",
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

		begin, err := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", beginDateVar))
		if err != nil {
			log.Fatalf("[F] Parse begin date failure, nest error: %v, begin: %s", err, beginDateVar)
		}
		end, err := time.Parse(time.DateTime, fmt.Sprintf("%s 23:59:59", endDateVar))
		if err != nil {
			log.Fatalf("[F] Parse end date failure, nest error: %v, end: %s", err, endDateVar)
		}

		stub, closeFunc, err := grpcclient.NewStorageWithEtcd()
		if err != nil {
			log.Fatalf("[F] New storage client failure, nest error: %v", err)
		}
		defer closeFunc()

		for begin.Before(end) {
			var (
				date = begin.Format(time.DateOnly)
				now  = time.Now()
			)
			resp, err := stub.ArchiveMetadata(context.Background(), &wrapperspb.StringValue{Value: date})
			if err != nil {
				log.Printf("[E] Archive data failure, nest error: %v, date: %v", err, date)
			} else {
				fmt.Printf("[I] Archive data success, date: %v, affected-stock: %v, affected-day: %v, affetced-week: %v, cost: %v\r\n", date, resp.AffectedStock, resp.AffectedQuoteDay, resp.AffectedQuoteWeek, time.Since(now))
			}
			begin = begin.Add(24 * time.Hour)

		}
	},
}
