package service

import (
	"context"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-cron"
	"github.com/eviltomorrow/king/lib/setting"
)

func Callback(schedulerId string, e error) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), setting.GRPC_UNARY_TIMEOUT_10_SECOND)
	defer cancel()

	code, msg := func() (pb.CallbackRequest_StatusCode, string) {
		if e != nil {
			return pb.CallbackRequest_FAILURE, e.Error()
		} else {
			return pb.CallbackRequest_SUCCESS, ""
		}
	}()
	if _, err := client.DefaultScheduler.Callback(ctx, &pb.CallbackRequest{SchedulerId: schedulerId, Code: code, ErrorMsg: msg}); err != nil {
		return "", err
	}

	return schedulerId, nil
}
