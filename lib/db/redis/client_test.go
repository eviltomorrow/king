package redis

import (
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	DSN = "redis://:admin123@localhost:6379/1?protocol=3"
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
}
