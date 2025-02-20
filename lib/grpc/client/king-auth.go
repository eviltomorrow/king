package client

import (
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/client/internal"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-auth"
)

var DefaultPassport pb.PassportClient

func InitPassport() error {
	client, shutdown, err := NewAuthWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultPassport = client
	return nil
}

func NewAuthWithEtcd() (pb.PassportClient, func() error, error) {
	target := "etcd:///grpclb/king-auth"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewPassportClient(conn), func() error { return conn.Close() }, nil
}

func NewAuthWithTarget(target string) (pb.PassportClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewPassportClient(conn), func() error { return conn.Close() }, nil
}
