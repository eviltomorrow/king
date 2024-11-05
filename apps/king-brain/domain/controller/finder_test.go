package controller

import (
	"context"
	"log"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDiscoverPossibleChance(t *testing.T) {
	instance := &Finder{}
	if _, err := instance.DiscoverPossibleChance(context.Background(), &wrapperspb.StringValue{
		Value: "2024-11-03",
	}); err != nil {
		log.Fatal(err)
	}
}
