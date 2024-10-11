package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var DB *mongo.Client

func InitMongoDB(c *Config) (func() error, error) {
	client, err := tryConnect(c)
	if err != nil {
		return nil, err
	}
	DB = client

	return func() error {
		if DB == nil {
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ConnectTimeout))
		defer cancel()

		return DB.Disconnect(ctx)
	}, nil
}

func tryConnect(c *Config) (*mongo.Client, error) {
	i := 1
	for {
		pool, err := buildMongoDB(c)
		if err == nil {
			return pool, nil
		}
		zlog.Error("connect to mongodb failure", zap.Error(err))
		i++

		if i > c.StartupRetryTimes {
			return nil, err
		}

		time.Sleep(time.Duration(c.StartupRetryPeriod))
	}
}

func buildMongoDB(c *Config) (*mongo.Client, error) {
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
