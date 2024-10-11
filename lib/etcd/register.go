package etcd

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/lib/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	ServicePrefix = "grpclb"
)

func RegisterService(ctx context.Context, service string, host string, port int, ttl int64) (func() error, error) {
	if Client == nil {
		return nil, fmt.Errorf("panic: etcd's client is nil")
	}

	leaseResp, err := Client.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	var leaseID = &leaseResp.ID

	key, value := fmt.Sprintf("/%s/%s/%s:%d", ServicePrefix, service, host, port), fmt.Sprintf("%s:%d", host, port)
	_, err = Client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}

	keepAlive, err := Client.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return nil, err
	}

	go func() {
	keep:
		for {
			select {
			case <-Client.Ctx().Done():
				return

			case k, ok := <-keepAlive:
				if !ok {
					break keep
				}
				_ = k

			case <-ctx.Done():
				return
			}
		}

	release:
		leaseResp, err := Client.Grant(ctx, ttl)
		if err != nil {
			zlog.Error("grant lease failure", zap.Error(err))
			goto release
		}

		key, value := fmt.Sprintf("/%s/%s:%d", service, host, port), fmt.Sprintf("%s:%d", host, port)
		_, err = Client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
		if err != nil {
			zlog.Error("put k/v failure", zap.Error(err), zap.String("key", key), zap.String("value", value))
			goto release
		}

		keepAlive, err = Client.KeepAlive(ctx, leaseResp.ID)
		if err != nil {
			zlog.Error("keepalive failure", zap.Error(err), zap.Any("leaseID", leaseResp.ID))
			goto release
		}
		leaseID = &leaseResp.ID

		goto keep
	}()
	revokeFunc := func() error {
		_, err = Client.Revoke(ctx, *leaseID)
		return err
	}

	return revokeFunc, nil
}
