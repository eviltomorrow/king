package metadata

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

var (
	inFlightSem = make(chan struct{}, 1)
)

func genCode(baseCodeList []string) chan string {
	data := make(chan string, 64)
	go func() {
		for _, code := range baseCodeList {
			result, err := genRangeCode(code)
			if err != nil {
				zlog.Error("build range code failure", zap.Error(err))
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
	if !strings.HasPrefix(baseCode, "sh") && !strings.HasPrefix(baseCode, "sz") && !strings.HasPrefix(baseCode, "bj") {
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
