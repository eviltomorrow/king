package conf

import (
	"bytes"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type NTFY struct {
	Scheme   string `json:"scheme"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
}

func (n *NTFY) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(n)
	return string(buf)
}

func LoadNTFY(path string) (*NTFY, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := bytes.TrimSpace(buf)
	s := &NTFY{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}
