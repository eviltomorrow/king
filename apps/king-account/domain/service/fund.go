package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"unsafe"

	"github.com/eviltomorrow/king/apps/king-account/domain/db"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/shopspring/decimal"
)

var FundLimitPerUserId int64 = 3

type AccountStatus int8

const (
	NORMAL AccountStatus = 0
	FREEZE
)

type Account struct {
	AliasName        string          `json:"alias_name"`
	UserId           string          `json:"user_id"`
	FundNo           string          `json:"fund_no"`
	OpeningCash      decimal.Decimal `json:"opening_cash"`
	EndCash          decimal.Decimal `json:"end_cash"`
	YesterdayEndCash decimal.Decimal `json:"yesterday_end_cash"`
	Status           int8            `json:"status"`
	InitDatetime     time.Time       `json:"init_datetime"`
}

func FundWithInitAccount(ctx context.Context, aliasName, userId string, openingCash decimal.Decimal) error {
	if aliasName == "" || userId == "" || openingCash.InexactFloat64() <= 0 {
		return fmt.Errorf("alias_name/user_id/opening_cash is nil")
	}

	count, err := db.FundWithCountByUserId(ctx, mysql.DB, userId)
	if err != nil {
		return err
	}
	if count >= FundLimitPerUserId {
		return fmt.Errorf("init account has reached the maximum: %d", FundLimitPerUserId)
	}

	cash := openingCash.InexactFloat64()
	fund := &db.Fund{
		AliasName:        aliasName,
		UserId:           userId,
		FundNo:           randomFundNo(),
		OpeningCash:      cash,
		EndCash:          cash,
		YesterdayEndCash: cash,
		Status:           int8(NORMAL),
		InitDatetime:     time.Now(),
	}

	_, err = db.FundWithInsertOne(ctx, mysql.DB, fund)
	if err != nil {
		return err
	}

	return err
}

func FundWithModifyAliasName(ctx context.Context, aliasName string, fundNo string) error {
	if aliasName == "" {
		return fmt.Errorf("alias_name is nil")
	}
	if fundNo == "" {
		return fmt.Errorf("fund_no is nil")
	}

	_, err := db.FundWithUpdateAliasName(ctx, mysql.DB, aliasName, fundNo)
	return err
}

func FundWithListByUserId(ctx context.Context, userId string) ([]*Account, error) {
	if userId == "" {
		return nil, fmt.Errorf("user_id is nil")
	}

	funds, err := db.FundWithSelectManyByUserId(ctx, mysql.DB, userId)
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0, len(funds))
	for _, fund := range funds {
		accounts = append(accounts, &Account{
			AliasName:        fund.AliasName,
			UserId:           fund.UserId,
			FundNo:           fund.FundNo,
			OpeningCash:      decimal.NewFromFloat(fund.OpeningCash),
			EndCash:          decimal.NewFromFloat(fund.EndCash),
			YesterdayEndCash: decimal.NewFromFloat(fund.YesterdayEndCash),
			Status:           fund.Status,
			InitDatetime:     fund.InitDatetime,
		})
	}
	return accounts, nil
}

func FundWithOneByFundNo(ctx context.Context, fundNo string) (*Account, error) {
	if fundNo == "" {
		return nil, fmt.Errorf("fund_no is nil")
	}

	fund, err := db.FundWithSelectOne(ctx, mysql.DB, fundNo)
	if err != nil {
		return nil, err
	}

	account := &Account{
		AliasName:        fund.AliasName,
		UserId:           fund.UserId,
		FundNo:           fund.FundNo,
		OpeningCash:      decimal.NewFromFloat(fund.OpeningCash),
		EndCash:          decimal.NewFromFloat(fund.EndCash),
		YesterdayEndCash: decimal.NewFromFloat(fund.YesterdayEndCash),
		Status:           fund.Status,
		InitDatetime:     fund.InitDatetime,
	}

	return account, nil
}

const letters = "1234567890"

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randomFundNo() string {
	return fmt.Sprintf("3%s", randomNumString(11))
}

func randomNumString(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())

	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
