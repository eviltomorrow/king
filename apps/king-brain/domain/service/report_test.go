package service

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestReport(t *testing.T) {
	d := time.Date(2024, time.November, 10, 0, 0, 0, 0, time.Local)
	status, err := ReportMarketStatus(context.Background(), d, "day")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)
}
