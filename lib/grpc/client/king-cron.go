package client

import (
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/client/internal"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-cron"
)

var DefaultScheduler pb.SchedulerClient

func InitScheduler() error {
	client, shutdown, err := NewSchedulerWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultScheduler = client
	return nil
}

func NewSchedulerWithEtcd() (pb.SchedulerClient, func() error, error) {
	target := "etcd:///grpclb/king-cron"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewSchedulerClient(conn), func() error { return conn.Close() }, nil
}

func NewSchedulerWithTarget(target string) (pb.SchedulerClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewSchedulerClient(conn), func() error { return conn.Close() }, nil
}
