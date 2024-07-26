package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/snowflake"
	"github.com/stretchr/testify/assert"
)

var fund1 = &Fund{
	FundNo:       "",
	UserId:       "2",
	OpeningCash:  600_000_000,
	EndCash:      600_000_000,
	Status:       0,
	InitDatetime: time.Date(2024, time.January, 10, 10, 0o0, 0o0, 0o0, time.Local),
}

var fund2 = &Fund{
	FundNo:       "",
	UserId:       "2",
	OpeningCash:  50_000_000,
	EndCash:      50_000_000,
	Status:       0,
	InitDatetime: time.Date(2024, time.January, 10, 10, 0o0, 0o0, 0o0, time.Local),
}

var fund3 = &Fund{
	FundNo:       "",
	UserId:       "2",
	OpeningCash:  300_000,
	EndCash:      300_000,
	Status:       1,
	InitDatetime: time.Date(2024, time.January, 10, 10, 0o0, 0o0, 0o0, time.Local),
}

func TruncateFund() error {
	_, err := mysql.DB.Exec("truncate table fund")
	return err
}

func InitFund() {
	TruncateFund()
}

func TestFundWithInsertOne(t *testing.T) {
	_assert := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitFund()

	for _, fund := range []*Fund{fund1, fund2, fund3} {
		fund.FundNo = snowflake.GenerateID()
		affected, err := FundWithInsertOne(ctx, mysql.DB, fund)
		_assert.Nil(err)
		_assert.Equal(affected, int64(1))

		fundInserted, err := FundWithSelectOne(ctx, mysql.DB, fund.FundNo)
		_assert.Nil(err)
		_assert.Equal(fundInserted.FundNo, fund.FundNo)
		_assert.Equal(fundInserted.UserId, fund.UserId)
		_assert.Equal(fundInserted.OpeningCash, fund.OpeningCash)
		_assert.Equal(fundInserted.EndCash, fund.EndCash)
		_assert.Equal(fundInserted.Status, fund.Status)
		_assert.Equal(fundInserted.InitDatetime, fund.InitDatetime)
	}

	affected, err := FundWithInsertOne(ctx, mysql.DB, nil)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)
}

func TestFundWithDeleteOne(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitFund()

	for _, fund := range []*Fund{fund1, fund2, fund3} {
		fund.FundNo = snowflake.GenerateID()
		affected, err := FundWithInsertOne(ctx, mysql.DB, fund)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)

		affected, err = FundWithDeleteOne(ctx, mysql.DB, fund.FundNo)
		_assert.Nil(err)
		_assert.Equal(int64(1), affected)

		_, err = FundWithSelectOne(ctx, mysql.DB, fund.FundNo)
		_assert.NotNil(err)
		_assert.Equal(sql.ErrNoRows, err)
	}

	id := snowflake.GenerateID()
	affected, err := FundWithDeleteOne(ctx, mysql.DB, id)
	_assert.Nil(err)
	_assert.Equal(int64(0), affected)
}

func TestFundWithUpdateOne(t *testing.T) {
	_assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	InitFund()

	fund1.FundNo = snowflake.GenerateID()
	affected, err := FundWithInsertOne(ctx, mysql.DB, fund1)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	fund1.Status = 2
	fund1.EndCash = 300000
	affected, err = FundWithUpdateOne(ctx, mysql.DB, fund1, fund1.FundNo, fund1.Version)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	fundUpdated, err := FundWithSelectOne(ctx, mysql.DB, fund1.FundNo)
	_assert.Nil(err)
	_assert.Equal(int8(2), fundUpdated.Status)
	_assert.Equal(fund1.EndCash, fundUpdated.EndCash)

	affected, err = FundWithUpdateOne(ctx, mysql.DB, nil, fund1.FundNo, fund1.Version)
	_assert.NotNil(err)
	_assert.Equal(int64(0), affected)
}
