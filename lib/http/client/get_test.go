package client

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	resp, err := Get("http://www.baidfsba13du1.com", 10*time.Second, defaultHttpHeader, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
