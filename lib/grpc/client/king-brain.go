package client

import (
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/client/internal"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-brain"
)

var DefaultFinder pb.FinderClient

func InitBrain() error {
	client, shutdown, err := NewBrainWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultFinder = client
	return nil
}

func NewBrainWithEtcd() (pb.FinderClient, func() error, error) {
	target := "etcd:///grpclb/king-brain"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewFinderClient(conn), func() error { return conn.Close() }, nil
}

func NewBrainWithTarget(target string) (pb.FinderClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewFinderClient(conn), func() error { return conn.Close() }, nil
}
