package conf

import (
	"bytes"
	"encoding/json"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type NFTY struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (n *NFTY) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(n)
	return string(buf)
}

func FindNTFY(path string) (*NFTY, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := bytes.TrimSpace(buf)
	s := &NFTY{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}
