package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/infrastructure"
	"github.com/eviltomorrow/king/lib/timeutil"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type client struct {
	ConnectTimeout     timeutil.Duration `json:"connect_timeout"`
	StartupRetryTimes  int               `json:"startup_retry_times"`
	StartupRetryPeriod timeutil.Duration `json:"startup_retry_period"`

	Endpoints []string `json:"endpoints"`
}

func (c *client) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

var Client *clientv3.Client

func init() {
	infrastructure.Register("etcd", &client{
		ConnectTimeout:     timeutil.Duration(3 * time.Second),
		StartupRetryTimes:  3,
		StartupRetryPeriod: timeutil.Duration(10 * time.Second),
	})
}

func (c *client) Connect() error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: time.Duration(c.ConnectTimeout),
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
		return err
	}

	var (
		i = 1
	)
	for {
		err := c.statusClient(client)
		if err == nil {
			break
		}
		zlog.Error("connect to etcd faialure", zap.Error(err))
		i++

		if i > c.StartupRetryTimes {
			return err
		}

		time.Sleep(time.Duration(c.StartupRetryPeriod))
	}
	Client = client

	return nil
}

func (c *client) Close() error {
	if Client == nil {
		return nil
	}

	Client.Sync(context.Background())
	return Client.Close()
}

func (c *client) UnMarshalConfig(config []byte) error {
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(config, c); err != nil {
		return err
	}
	return nil
}

func (c *client) statusClient(client *clientv3.Client) error {
	for _, endpoint := range c.Endpoints {
		if err := c.statusEndpoint(client, endpoint); err == nil {
			return nil
		}
	}
	return fmt.Errorf("connect to etcd service failure, nest endpoints: %v", c.Endpoints)
}

func (c *client) statusEndpoint(client *clientv3.Client, endpoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ConnectTimeout))
	defer cancel()

	_, err := client.Status(ctx, endpoint)
	return err
}
