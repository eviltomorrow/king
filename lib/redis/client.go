package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/infrastructure"
	"github.com/eviltomorrow/king/lib/timeutil"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type client struct {
	ConnectTimeout     timeutil.Duration `json:"connect_timeout"`
	StartupRetryTimes  int               `json:"startup_retry_times"`
	StartupRetryPeriod timeutil.Duration `json:"startup_retry_period"`

	DSN string `json:"dsn"`
}

func (c *client) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

var Client *redis.Client

func init() {
	infrastructure.Register("redis", &client{
		ConnectTimeout:     timeutil.Duration(3 * time.Second),
		StartupRetryTimes:  3,
		StartupRetryPeriod: timeutil.Duration(10 * time.Second),
	})
}

func (c *client) Connect() error {
	var (
		pool *redis.Client
		err  error

		i = 1
	)
	for {
		pool, err = c.buildRedis()
		if err == nil {
			break
		}
		zlog.Error("connect to redis faialure", zap.Error(err))
		i++

		if i > c.StartupRetryTimes {
			return err
		}

		time.Sleep(time.Duration(c.StartupRetryPeriod))
	}
	Client = pool

	return nil
}

func (c *client) Close() error {
	if Client == nil {
		return nil
	}

	return Client.Close()
}

func (c *client) UnMarshalConfig(config []byte) error {
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(config, c); err != nil {
		return err
	}
	return nil
}

func (c *client) buildRedis() (*redis.Client, error) {
	if c.DSN == "" {
		return nil, fmt.Errorf("redis: no DSN set")
	}

	opts, err := redis.ParseURL(c.DSN)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ConnectTimeout))
	defer cancel()

	status := client.Ping(ctx)
	if err := status.Err(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
