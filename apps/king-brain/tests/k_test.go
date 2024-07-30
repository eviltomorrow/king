package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/king/apps/king-brain/domain/chart"
	"github.com/eviltomorrow/king/apps/king-brain/domain/data"
	pb "github.com/eviltomorrow/king/lib/grpc/pb/king-storage"
)

func TestNewK(t *testing.T) {
	quotes, err := data.FetchQuote(context.Background(), time.Now(), "sh600519")
	if err != nil {
		log.Fatal(err)
	}

	k, err := chart.NewK(context.Background(), &pb.Stock{Name: "贵州茅台", Code: "sh600519"}, quotes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k)
}
