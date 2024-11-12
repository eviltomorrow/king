package metadata

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
	"github.com/eviltomorrow/king/lib/grpc/transformer"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var ReportCommand = &cobra.Command{
	Use:   "report",
	Short: "报告数据",
	Run: func(cmd *cobra.Command, args []string) {
		date, err := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", begin))
		if err != nil {
			log.Printf("[F] 转换日期失败, nest error: %v, begin: %s", err, begin)
			return
		}
		if begin == "" {
			date = time.Now()
		}
		if err := report(context.Background(), "daily", date.Format(time.DateOnly)); err != nil {
			log.Printf("[F] Report failure, nest error: %v", err)
		} else {
			log.Printf("=> 日期：%s，报告已生成，请查看邮箱", begin)
		}
	},
}

func init() {
	ReportCommand.PersistentFlags().StringVar(&begin, "date", "", "指定日期")
	ReportCommand.PersistentFlags().StringVar(&mode, "mode", "daily", "指定类型")

	ReportCommand.PersistentFlags().StringVar(&IP, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func report(ctx context.Context, mode, date string) error {
	clientBrain, closeFunc, err := client.NewBrainWithTarget(fmt.Sprintf("%s:50005", IP))
	if err != nil {
		return err
	}
	defer closeFunc()

	resp, err := clientBrain.ReportDaily(ctx, &wrapperspb.StringValue{Value: date})
	if err != nil {
		return err
	}

	value := make(map[string]string)
	data := transformer.GenerateMarketStatusToMap(resp)
	for k, v := range data {
		value[k] = fmt.Sprintf("%v", v)
	}
	clientTamplate, closeFunc, err := client.NewTemplateWithTarget(fmt.Sprintf("%s:50002", IP))
	if err != nil {
		return err
	}
	defer closeFunc()

	text, err := clientTamplate.Render(context.Background(), &pb.RenderRequest{
		TemplateName: fmt.Sprintf("%s-report.html", mode),
		Data:         value,
	})
	if err != nil {
		return err
	}

	clientEmail, closeFunc, err := client.NewEmailWithTarget(fmt.Sprintf("%s:50002", IP))
	if err != nil {
		return err
	}
	defer closeFunc()

	_, err = clientEmail.Send(context.Background(), &pb.Mail{
		To: []*pb.Contact{
			{Name: "liarsa", Address: "eviltomorrow@163.com"},
		},
		Subject: fmt.Sprintf("简报(%s)", date),
		Body:    text.Value,
	})
	if err != nil {
		return err
	}

	return err
}
