package service

import (
	"context"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-cron"
	"github.com/eviltomorrow/king/lib/setting"
)

func Callback(schedulerId string, e error) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10SECOND)
	defer cancel()

	code := pb.CallbackRequest_SUCCESS
	msg := ""

	if e != nil {
		code = pb.CallbackRequest_FAILURE
		msg = e.Error()
	}
	if _, err := client.DefalutScheduler.Callback(ctx, &pb.CallbackRequest{SchedulerId: schedulerId, Code: code, ErrorMsg: msg}); err != nil {
		return "", err
	}

	return schedulerId, nil
}
