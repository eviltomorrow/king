package datasource

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/eviltomorrow/king/lib/http/client"
	"github.com/eviltomorrow/king/lib/model"
	"github.com/eviltomorrow/king/lib/setting"
	"github.com/eviltomorrow/king/lib/zlog"
	"go.uber.org/zap"
)

const (
	suspendNormal    = "正常"
	suspendOneHour   = "停牌一小时"
	suspendOneDay    = "停牌一天"
	suspendKeep      = "连续停牌"
	suspendMid       = "盘中停牌"
	suspendHalfOfDay = "停牌半天"
	suspendPause     = "暂停"
	suspendNoRecord  = "无该记录"
	suspendUnlisted  = "未上市"
	suspendDelist    = "退市"
	suspendUnknown   = "未知"
)

var (
	SinaHeader = map[string]string{
		"Referer":                   "https://finance.sina.com.cn",
		"Connection":                "keep-alive",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Encoding":           "gzip, deflate",
		"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,da;q=0.7,pt;q=0.6,ja;q=0.5",
	}
	SinaMatcher = map[string][]int{
		"sh68":  {33, 34},
		"sh60":  {33},
		"sz0":   {33},
		"sz3":   {33, 34},
		"sh0":   {33, 34},
		"bj8":   {38, 39},
		"bj920": {38, 39},
	}
)

func FetchMetadataFromSina(codes []string) ([]*model.Metadata, error) {
	if len(codes) == 0 {
		return nil, fmt.Errorf("codes is nil")
	}
	url := fmt.Sprintf("https://hq.sinajs.cn/list=%s", strings.Join(codes, ","))

	data, err := client.Get(url, setting.DEFUALT_HANDLE_10_SECOND, SinaHeader, nil)
	if err != nil {
		return nil, fmt.Errorf("url: %v, nest error: %v", url, err)
	}

	zlog.Debug("Get stock data", zap.String("data", data))
	var (
		result = make([]*model.Metadata, 0, len(codes))
		exists = make(map[string]struct{}, len(codes))
	)
	kv, err := parseSinaDataToMap(data)
	if err != nil {
		zlog.Error("parseSinaDataToMap failure", zap.Error(err), zap.String("data", data))
	}
	for key, val := range kv {
		metadata, err := parseSinaLineToMetadata(key, val)
		if err != nil {
			zlog.Error("parseSinaLineToMetadata failure", zap.Error(err), zap.String("key", key), zap.String("val", val))
		}
		if metadata != nil {
			if _, ok := exists[metadata.Code]; !ok {
				result = append(result, metadata)
				exists[metadata.Code] = struct{}{}
			}
		}
	}
	return result, nil
}

func parseSinaDataToMap(data string) (map[string]string, error) {
	result := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		if !strings.HasPrefix(text, "var") || !strings.HasSuffix(text, ";") {
			return nil, fmt.Errorf("invalid line data")
		}

		n := strings.Index(text, "=")
		if n == -1 {
			return nil, fmt.Errorf("invalid line data")
		}

		code := strings.Replace(text[:n], "var hq_str_", "", -1)
		result[code] = text
	}
	return result, nil
}

func parseSinaLineToMetadata(code, data string) (*model.Metadata, error) {
	if data == "" {
		return nil, nil
	}
	var (
		begin = strings.Index(strings.TrimSpace(data), `"`)
		end   = strings.LastIndex(strings.TrimSpace(data), `"`)
	)

	if begin == -1 || end == -1 || begin == end {
		return nil, fmt.Errorf("panic: begin: %v, end: %v", begin, end)
	} else {
		if strings.TrimSpace(data[begin+1:end]) == "" {
			return nil, nil
		}
	}
	attr := strings.Split(data[begin+1:end], ",")
	if len(attr) == 1 {
		return nil, fmt.Errorf("panic: attr format is unknown, nest attr: %v", attr)
	}

	if len(attr) >= 2 && attr[len(attr)-1] == "" {
		attr = attr[:len(attr)-1]
	}

	var support bool
loop:
	for key, val := range SinaMatcher {
		if strings.HasPrefix(code, key) {
			for _, v := range val {
				if len(attr) == v {
					support = true
					break loop
				}
			}
			return nil, fmt.Errorf("format is changed[%v], expect: %v, actual: %v", code, val, len(attr))
		}
	}
	if !support {
		return nil, fmt.Errorf("not support code[%v]", code)
	}
	// switch {
	// case strings.HasPrefix(code, "sh68"):
	// 	if len(attr) != matcher["sh68"] {
	// 		return nil, fmt.Errorf("format is changed[sh68xxxx], expect: %v, actual: %v", matcher["sh68"], len(attr))
	// 	}
	// case strings.HasPrefix(code, "sh60"):
	// 	if len(attr) != matcher["sh60"] {
	// 		return nil, fmt.Errorf("format is changed[sh60xxxx] expect: %v, actual: %v", matcher["sh60"], len(attr))
	// 	}
	// case strings.HasPrefix(code, "sz0"):
	// 	if len(attr) != matcher["sz0"] {
	// 		return nil, fmt.Errorf("format is changed[sz0xxxxx] expect: %v, actual: %v", matcher["sz0"], len(attr))
	// 	}
	// case strings.HasPrefix(code, "sz3"):
	// 	if len(attr) != matcher["sz3"] {
	// 		return nil, fmt.Errorf("format is changed[sz3xxxxx] expect: %v, actual: %v", matcher["sz3"], len(attr))
	// 	}
	// default:
	// 	return nil, fmt.Errorf("panic: no support code[%v]", code)
	// }

	md := &model.Metadata{
		Code: code,
	}
	for i, val := range attr {
		switch i {
		case 0:
			md.Name = val
		case 1:
			md.Open = atof64(md.Name, i, val)
		case 2:
			md.YesterdayClosed = atof64(md.Name, i, val)
		case 3:
			md.Latest = atof64(md.Name, i, val)
		case 4:
			md.High = atof64(md.Name, i, val)
		case 5:
			md.Low = atof64(md.Name, i, val)
		case 8:
			md.Volume = atoi64(md.Name, i, val)
		case 9:
			md.Account = atof64(md.Name, i, val)
		case 30:
			md.Date = val
		case 31:
			md.Time = val
		case 32:
			md.Suspend = getSuspendDesc(val)
		default:
		}
	}
	return md, nil
}

// getSuspendDesc get suspend desc
func getSuspendDesc(val string) string {
	switch {
	case val == "00":
		return suspendNormal
	case val == "01":
		return suspendOneHour
	case val == "02":
		return suspendOneDay
	case val == "03":
		return suspendKeep
	case val == "04":
		return suspendMid
	case val == "05":
		return suspendHalfOfDay
	case val == "07":
		return suspendPause
	case val == "-1":
		return suspendNoRecord
	case val == "-2":
		return suspendUnlisted
	case val == "-3":
		return suspendDelist
	default:
		return suspendUnknown
	}
}
