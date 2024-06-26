package service

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/eviltomorrow/king/apps/king-collector/domain/service/datasource"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

var (
	inFlightSem = make(chan struct{}, 1)
	CodeList    = []string{
		"sh688***",
		"sh605***",
		"sh603***",
		"sh601***",
		"sh600***",
		"sz300***",
		"sz0030**",
		"sz002***",
		"sz001***",
		"sz000***",
	}
	FetchFactories = map[string]func([]string) ([]*model.Metadata, error){
		"sina":   datasource.FetchMetadataFromSina,
		"net126": datasource.FetchMetadataFromNet126,
	}
	Size         = 30
	Limit        = 3
	Timeout      = 10 * time.Second
	RandomPeriod = [2]int{20, 60}
)

func genCode() chan string {
	data := make(chan string, 64)
	go func() {
		for _, code := range CodeList {
			result, err := genRangeCode(code)
			if err != nil {
				zlog.Error("Build range code failure", zap.Error(err))
				continue
			}
			for _, r := range result {
				data <- r
			}
		}
		close(data)
	}()
	return data
}

func genRangeCode(baseCode string) ([]string, error) {
	if len(baseCode) != 8 {
		return nil, fmt.Errorf("code length must be 8, code is [%s]", baseCode)
	}
	if !strings.HasPrefix(baseCode, "sh") && !strings.HasPrefix(baseCode, "sz") {
		return nil, fmt.Errorf("code must be start with [sh/sz], code is [%s]", baseCode)
	}

	if !strings.Contains(baseCode, "*") {
		return []string{baseCode}, nil
	}

	var (
		n      = strings.Index(baseCode, "*")
		prefix = baseCode[:n]
		codes  = make([]string, 0, int(math.Pow10(8-n)))
	)

	var builder strings.Builder
	builder.Grow(8)

	next := int(math.Pow10(8-n)) - 1
	mid := ""
	count := -1
	changed := false
	for i := next; i >= 0; i-- {
		if i == next && i != 0 {
			next = i / 10
			count++
			changed = true
			mid = ""
		} else {
			changed = false
		}

		if changed {
			for j := 0; j < count; j++ {
				mid += "0"
			}
		}

		builder.WriteString(prefix)
		builder.WriteString(mid)
		builder.WriteString(strconv.Itoa(i))
		codes = append(codes, builder.String())
		builder.Reset()
	}
	return codes, nil
}
