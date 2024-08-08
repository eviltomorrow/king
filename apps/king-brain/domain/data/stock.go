package data

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type Stock struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Suspend string `json:"suspend"`
}

func (s *Stock) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s)
	if err != nil {
		return fmt.Sprintf("marshal stock failure, nest error: %v", err)
	}
	return string(buf)
}
