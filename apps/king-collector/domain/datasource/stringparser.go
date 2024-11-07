package datasource

import (
	"strconv"

	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

func atof64(name string, loc int, val string) float64 {
	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		zlog.Error("parse float64 failure", zap.String("name", name), zap.Int("loc", loc), zap.String("val", val))
		return 0
	}
	return f64
}

func atou64(name string, loc int, val string) uint64 {
	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		zlog.Error("parse uint64 failure", zap.String("name", name), zap.Int("loc", loc), zap.String("val", val))
		return 0
	}
	return u64
}

func atoi64(name string, loc int, val string) int64 {
	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		zlog.Error("parse int64 failure", zap.String("name", name), zap.Int("loc", loc), zap.String("val", val))
		return 0
	}
	return i64
}
