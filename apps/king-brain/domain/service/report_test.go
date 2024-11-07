package service

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestReport(t *testing.T) {
	status, err := ReportLatest(context.Background(), time.Now(), "day")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)
}
