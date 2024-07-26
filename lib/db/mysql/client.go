package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/infrastructure"
	"github.com/eviltomorrow/king/lib/timeutil"
	"github.com/eviltomorrow/king/lib/zlog"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type client struct {
	ConnectTimeout     timeutil.Duration `json:"connect_timeout"`
	StartupRetryTimes  int               `json:"startup_retry_times"`
	StartupRetryPeriod timeutil.Duration `json:"startup_retry_period"`

	DSN         string            `json:"dsn"`
	MinOpen     int               `json:"min_open"`
	MaxOpen     int               `json:"max_open"`
	MaxLifetime timeutil.Duration `json:"max_lifetime"`
}

func (c *client) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}

var DB *sql.DB

func init() {
	infrastructure.Register("mysql", &client{
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
		pool *sql.DB
		err  error

		i = 1
	)
	for {
		pool, err = c.buildMySQL()
		if err == nil {
			break
		}
		zlog.Error("connect to mysql faialure", zap.Error(err))
		i++
		if i > c.StartupRetryTimes {
			return err
		}

		time.Sleep(time.Duration(c.StartupRetryPeriod))
	}
	DB = pool

	return nil
}

func (c *client) Close() error {
	if DB == nil {
		return nil
	}

	return DB.Close()
}

func (c *client) UnMarshalConfig(config []byte) error {
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(config, c); err != nil {
		return err
	}
	return nil
}

func (c *client) buildMySQL() (*sql.DB, error) {
	if c.DSN == "" {
		return nil, fmt.Errorf("MySQL: no DSN set")
	}
	pool, err := sql.Open("mysql", c.DSN)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(time.Duration(c.MaxLifetime))
	pool.SetMaxOpenConns(c.MaxOpen)
	pool.SetMaxIdleConns(c.MinOpen)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ConnectTimeout))
	defer cancel()

	if err = pool.PingContext(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

// Exec exec mysql
type Exec interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
