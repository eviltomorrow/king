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
		var (
			begin = time.Now()
			err   error
		)
		if beginVar != "" {
			begin, err = time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:01", beginVar))
			if err != nil {
				log.Printf("[F] 转换日期失败, nest error: %v, begin: %s", err, beginVar)
				return
			}
		}

		if err := report(context.Background(), "daily", begin.Format(time.DateOnly)); err != nil {
			log.Printf("[F] Report failure, nest error: %v", err)
		} else {
			log.Printf("=> 日期：%s，报告已生成，请查看 ntfy 通知 app.", begin.Format(time.DateOnly))
		}
	},
}

func init() {
	ReportCommand.PersistentFlags().StringVar(&beginVar, "date", "", "指定日期")
	ReportCommand.PersistentFlags().StringVar(&modeVar, "mode", "daily", "指定类型")

	ReportCommand.PersistentFlags().StringVar(&IPVar, "ip", "127.0.0.1", "指定服务端 IP 地址")
}

func report(ctx context.Context, mode, date string) error {
	clientBrain, closeFunc, err := client.NewBrainWithTarget(fmt.Sprintf("%s:50005", IPVar))
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
	clientTamplate, closeFunc, err := client.NewTemplateWithTarget(fmt.Sprintf("%s:50002", IPVar))
	if err != nil {
		return err
	}
	defer closeFunc()

	text, err := clientTamplate.Render(context.Background(), &pb.RenderRequest{
		TemplateName: fmt.Sprintf("%s-report.txt", mode),
		Data:         value,
	})
	if err != nil {
		return err
	}

	clientNtfy, closeFunc, err := client.NewNtfyWithTarget(fmt.Sprintf("%s:50002", IPVar))
	if err != nil {
		return err
	}
	defer closeFunc()

	if _, err := clientNtfy.Send(context.Background(), &pb.Msg{
		Topic:    "SrxOPwCBiRWZUOq0",
		Message:  text.Value,
		Title:    fmt.Sprintf("%s 日 汇总", date),
		Priority: 3,
		Tags:     []string{"简报", "股票", "统计"},
	}); err != nil {
		return err
	}

	return err
}
