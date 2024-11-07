package data

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type Quote struct {
	Code            string  `json:"code"`
	Open            float64 `json:"opend"`
	Close           float64 `json:"close"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	YesterdayClosed float64 `json:"yesterday_closed"`
	Volume          int64   `json:"volume"`
	Account         float64 `json:"account"`
	Date            string  `json:"date"`
	NumOfYear       int32   `json:"num_of_year"`
}

func (q *Quote) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	if err != nil {
		return fmt.Sprintf("marshal quote failure, nest error: %v", err)
	}
	return string(buf)
}
