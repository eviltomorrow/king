package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/infrastructure"
	"github.com/eviltomorrow/king/lib/timeutil"
	"github.com/eviltomorrow/king/lib/zlog"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type client struct {
	ConnectTimeout     timeutil.Duration `json:"connect_timeout"`
	StartupRetryTimes  int               `json:"startup_retry_times"`
	StartupRetryPeriod timeutil.Duration `json:"startup_retry_period"`

	DSN         string            `json:"dsn"`
	MinOpen     uint64            `json:"min_open"`
	MaxOpen     uint64            `json:"max_open"`
	MaxLifetime timeutil.Duration `json:"max_lifetime"`
}

func (c *client) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

var DB *mongo.Client

func init() {
	infrastructure.Register("mongodb", &client{
		ConnectTimeout:     timeutil.Duration(3 * time.Second),
		StartupRetryTimes:  3,
		StartupRetryPeriod: timeutil.Duration(10 * time.Second),

		MinOpen:     5,
		MaxOpen:     10,
		MaxLifetime: timeutil.Duration(3 * time.Minute),
	})
}

func (c *client) Connect() error {
	var (
		pool *mongo.Client
		err  error

		i = 1
	)
	for {
		pool, err = c.buildMongoDB()
		if err == nil {
			break
		}
		zlog.Error("connect to mongodb faialure", zap.Error(err))
		i++

		if i > c.StartupRetryTimes {
			return err
		}

		time.Sleep(time.Duration(c.StartupRetryPeriod))
	}
	DB = pool

	return err
}

func (c *client) Close() error {
	if DB == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ConnectTimeout))
	defer cancel()

	return DB.Disconnect(ctx)
}

func (c *client) UnMarshalConfig(config []byte) error {
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(config, c); err != nil {
		return err
	}
	return nil
}

func (c *client) buildMongoDB() (*mongo.Client, error) {
	if c.DSN == "" {
		return nil, fmt.Errorf("MongoDB: no DSN set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ConnectTimeout))
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(c.DSN).SetMaxPoolSize(c.MaxOpen).SetMinPoolSize(c.MinOpen).SetMaxConnIdleTime(time.Duration(c.MaxLifetime)),
	)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
