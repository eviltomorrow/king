package controller

import (
	"context"
	"log"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFindPossibleChance(t *testing.T) {
	instance := &Finder{}
	if _, err := instance.FindPossibleChance(context.Background(), &wrapperspb.StringValue{
		Value: "2024-12-22",
	}); err != nil {
		log.Fatal(err)
	}
}
