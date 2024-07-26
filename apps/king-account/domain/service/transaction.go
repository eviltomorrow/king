package service

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type (
	SecuritiesType        int8
	TransactionRecordType int8
)

const (
	STOCK SecuritiesType = iota
	ETF
)

const (
	OPENING TransactionRecordType = iota
	BUY
	SELL
	CLEAR
)

type Securities struct {
	UserId string
	FundNo string

	Type       SecuritiesType
	Code       string
	Name       string
	Volumn     int64
	ClosePrice float64
}

func Buy(ctx context.Context, securities *Securities) error {
	if securities == nil {
		return fmt.Errorf("securities is nil")
	}

	fund, err := FundWithOneByFundNo(ctx, securities.FundNo)
	if err != nil {
		return err
	}

	stock, err := StockWithGetOne(ctx, securities.Code)
	if err != nil {
		return err
	}

	requiredFund := decimal.NewFromFloat(securities.ClosePrice).Mul(decimal.NewFromInt(securities.Volumn))
	if requiredFund.GreaterThan(fund.EndCash) {
		return fmt.Errorf("no enough fund")
	}

	_ = stock

	return nil
}

func Sell(ctx context.Context, securities *Securities) error {
	return nil
}
