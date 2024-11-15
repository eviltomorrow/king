package service

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestReport(t *testing.T) {
	d := time.Date(2024, time.November, 15, 0, 0, 0, 0, time.Local)
	status, err := ReportMarketStatus(context.Background(), d, "week")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)
}
