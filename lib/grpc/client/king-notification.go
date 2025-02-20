package client

import (
	"github.com/eviltomorrow/king/lib/finalizer"
	"github.com/eviltomorrow/king/lib/grpc/client/internal"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-notification"
)

var (
	DefaultEmail    pb.EmailClient
	DefaultNTFY     pb.NtfyClient
	DefaultTemplate pb.TemplateClient
)

func InitEmail() error {
	client, shutdown, err := NewEmailWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultEmail = client
	return nil
}

func InitNTFY() error {
	client, shutdown, err := NewNtfyWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultNTFY = client
	return nil
}

func InitTemplate() error {
	client, shutdown, err := NewTemplateWithEtcd()
	if err != nil {
		return err
	}
	finalizer.RegisterCleanupFuncs(shutdown)

	DefaultTemplate = client
	return nil
}

func NewEmailWithEtcd() (pb.EmailClient, func() error, error) {
	target := "etcd:///grpclb/king-notification"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}

func NewEmailWithTarget(target string) (pb.EmailClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewEmailClient(conn), func() error { return conn.Close() }, nil
}

func NewNtfyWithEtcd() (pb.NtfyClient, func() error, error) {
	target := "etcd:///grpclb/king-notification"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewNtfyClient(conn), func() error { return conn.Close() }, nil
}

func NewNtfyWithTarget(target string) (pb.NtfyClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewNtfyClient(conn), func() error { return conn.Close() }, nil
}

func NewTemplateWithEtcd() (pb.TemplateClient, func() error, error) {
	target := "etcd:///grpclb/king-notification"
	conn, err := internal.DialWithEtcd(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewTemplateClient(conn), func() error { return conn.Close() }, nil
}

func NewTemplateWithTarget(target string) (pb.TemplateClient, func() error, error) {
	conn, err := internal.DialWithDirect(target)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewTemplateClient(conn), func() error { return conn.Close() }, nil
}
