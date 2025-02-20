package client

import (
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/client/internal"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
)

var DefaultStorage pb.StorageClient

func InitStorage() error {
	client, shutdown, err := NewStorageWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultStorage = client
	return nil
}

func NewStorageWithEtcd() (pb.StorageClient, func() error, error) {
	target := "etcd:///grpclb/king-storage"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewStorageClient(conn), func() error { return conn.Close() }, nil
}

func NewStorageWithTarget(target string) (pb.StorageClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewStorageClient(conn), func() error { return conn.Close() }, nil
}
