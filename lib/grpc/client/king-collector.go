package client

import (
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/client/internal"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-collector"
)

var DefaultCollector pb.CollectorClient

func InitCollector() error {
	client, shutdown, err := NewCollectorWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultCollector = client
	return nil
}

func NewCollectorWithEtcd() (pb.CollectorClient, func() error, error) {
	target := "etcd:///grpclb/king-collector"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewCollectorClient(conn), func() error { return conn.Close() }, nil
}

func NewCollectorWithTarget(target string) (pb.CollectorClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewCollectorClient(conn), func() error { return conn.Close() }, nil
}
