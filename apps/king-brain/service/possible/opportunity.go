package possible

import (
	"context"
	"fmt"

	"github.com/eviltomorrow/king/apps/king-brain/service/data"
)

type Opportunity struct {
	Score float64
}

type Stock struct {
	Code string
	Name string
}

func FindOpportunities(ctx context.Context, date string) ([]*Opportunity, error) {
	ch, err := data.NewDataWrapperChannel(data.DAY, date)
	if err != nil {
		return nil, err
	}

	for day := range ch {
		fmt.Println(day)
	}
	return nil, nil
}
