package middleware

import (
	"github.com/eviltomorrow/king/lib/zlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func InitLogger(c *zlog.Config) (func() error, error) {
	access, _, err := zlog.InitLogger(c)
	if err != nil {
		return nil, err
	}
	logger = access
	return logger.Sync, nil
}

func ServerLogInterceptor() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("",
				zap.Time("start-time", v.StartTime),
				zap.Int("status", v.Status),
				zap.String("method", v.Method),
				zap.String("remote-ip", v.RemoteIP),
				zap.String("uri", v.URI),
				zap.Duration("latency", v.Latency),
				zap.Error(v.Error),
			)

			return nil
		},
	})

}
