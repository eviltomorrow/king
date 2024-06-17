package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	DSN string
	RDB *redis.Client

	RetryTimes = 3
	Period     = 10 * time.Second
)

var DefaultConnectTimeout = 5 * time.Second

func Connect() error {
	var (
		pool *redis.Client
		err  error

		i = 1
	)
	for {
		if i > RetryTimes {
			if err != nil {
				return err
			}
			return fmt.Errorf("panic: connect Redis failure, err is nil?")
		}
		pool, err = buildRedis(DSN)
		if err == nil {
			break
		}

		log.Printf("[E] Try to connect to Redis, retry: %d, nest error: %v", i, err)
		i++
		time.Sleep(Period)
	}
	RDB = pool

	return nil
}

// Close close redis
func Close() error {
	if RDB == nil {
		return nil
	}

	return RDB.Close()
}

func buildRedis(dsn string) (*redis.Client, error) {
	if dsn == "" {
		return nil, fmt.Errorf("redis: no DSN set")
	}

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	status := client.Ping(ctx)
	if err := status.Err(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
