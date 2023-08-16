package etcd

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	Endpoints      []string
	ConnectTimeout = 5 * time.Second
	ExecuteTimeout = 20 * time.Second

	RetryTimes = 3
	Period     = 10 * time.Second
)

func NewClient() (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: ConnectTimeout,
		LogConfig: &zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.ErrorLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	})
	if err != nil {
		return nil, err
	}

	var (
		i = 1
		e error
	)
	for {
		if i > RetryTimes {
			if e != nil {
				return nil, e
			}
			return nil, fmt.Errorf("panic: connect etcd failure, err is nil?")
		}

		e = statusClient(client)
		if e == nil {
			break
		}
		if e != nil {
			log.Printf("[W] Connect to Etcd. Retry: %d, nest error: %v", i, e)
		}
		i++
		time.Sleep(Period)
	}
	return client, nil
}

func statusClient(client *clientv3.Client) error {
	for _, endpoint := range Endpoints {
		if err := statusEndpoint(client, endpoint); err == nil {
			return nil
		}
	}
	return fmt.Errorf("connect to etcd service failure, nest error: no valid endpoint, endpoints: %v", Endpoints)
}

func statusEndpoint(client *clientv3.Client, endpoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	_, err := client.Status(ctx, endpoint)
	return err
}
