package redis

import (
	"context"
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	DSN = "redis://:admin123@localhost:6379/0?protocol=3"
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	s := RDB.Set(context.Background(), "k", "v", 0)
	t.Log(s.Err())

	r := RDB.Get(context.Background(), "k")
	t.Log(r.Err())
	t.Log(r.Result())

	r = RDB.Get(context.Background(), "w")
	t.Log(r.Err())
	t.Log(r.Result())
}
