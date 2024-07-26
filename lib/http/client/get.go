package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Get(url string, timeout time.Duration, header map[string]string, body io.Reader) (string, error) {
	request, err := http.NewRequest("GET", url, body)
	if err != nil {
		return "", fmt.Errorf("new http request failure, nest error: %v", err)
	}

	for key, val := range header {
		request.Header.Add(key, val)
	}

	return Do(HTTPDefault, request, timeout)
}
