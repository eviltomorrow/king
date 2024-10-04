package callback

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/grpc/client"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-cron"
	"google.golang.org/grpc/metadata"
)

func Do(ctx context.Context, e error) (string, error) {
	params, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("context has no value")
	}

	data := params["scheduler_id"]
	if len(data) == 0 {
		return "", fmt.Errorf("context has no scheduler_id")
	}
	schedulerId := data[0]

	stub, shutdown, err := client.NewSchedulerWithEtcd()
	if err != nil {
		return "", err
	}
	defer shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	code := pb.CallbackRequest_SUCCESS
	msg := ""

	if e != nil {
		code = pb.CallbackRequest_FAILURE
		msg = e.Error()
	}
	if _, err = stub.Callback(ctx, &pb.CallbackRequest{SchedulerId: schedulerId, Code: code, ErrorMsg: msg}); err != nil {
		return "", err
	}

	return schedulerId, nil
}
