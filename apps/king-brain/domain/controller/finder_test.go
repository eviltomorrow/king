package controller

import (
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFindPossibleChance(t *testing.T) {
	instance := &Finder{}
	if _, err := instance.FindPossibleChance(context.Background(), &wrapperspb.StringValue{
		Value: time.Now().Format(time.DateOnly),
	}); err != nil {
		log.Fatal(err)
	}
}
