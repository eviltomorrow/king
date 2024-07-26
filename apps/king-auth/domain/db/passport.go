package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/orm"
)

func PassportWithSelectOneByAccount(ctx context.Context, exec mysql.Exec, account string) (*Passport, error) {
	passport := Passport{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&passport.Id,
			&passport.Account,
			&passport.Code,
			&passport.Salt,
			&passport.SaltPassword,
			&passport.Email,
			&passport.Phone,
			&passport.Status,
			&passport.CreateTimestamp,
			&passport.ModifyTimestamp,
		)
	}
	if err := orm.TableWithSelectOne(ctx, exec, TablePassportName, PassportFields, map[string]interface{}{FieldPassportAccount: account}, scan); err != nil {
		return nil, err
	}
	return &passport, nil
}

func PassportWithSelectOneById(ctx context.Context, exec mysql.Exec, id string) (*Passport, error) {
	passport := Passport{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&passport.Id,
			&passport.Account,
			&passport.Code,
			&passport.Salt,
			&passport.SaltPassword,
			&passport.Email,
			&passport.Phone,
			&passport.Status,
			&passport.CreateTimestamp,
			&passport.ModifyTimestamp,
		)
	}
	if err := orm.TableWithSelectOne(ctx, exec, TablePassportName, PassportFields, map[string]interface{}{FieldPassportId: id}, scan); err != nil {
		return nil, err
	}
	return &passport, nil
}

func PassportWithSelectRange(ctx context.Context, exec mysql.Exec, offset, limit int64) ([]*Passport, error) {
	passports := make([]*Passport, 0, limit)
	scan := func(rows *sql.Rows) error {
		for rows.Next() {
			passport := Passport{}
			if err := rows.Scan(
				&passport.Id,
				&passport.Account,
				&passport.Code,
				&passport.Salt,
				&passport.SaltPassword,
				&passport.Email,
				&passport.Phone,
				&passport.Status,
				&passport.CreateTimestamp,
				&passport.ModifyTimestamp,
			); err != nil {
				return err
			}
			passports = append(passports, &passport)
		}
		return nil
	}
	if err := orm.TableWithSelectRange(ctx, exec, TablePassportName, PassportFields, nil, nil, offset, limit, scan); err != nil {
		return nil, err
	}
	return passports, nil
}

func PassportWithUpdateStatus(ctx context.Context, exec mysql.Exec, status int8, id string) (int64, error) {
	if status < 0 {
		return 0, fmt.Errorf("status is invalid")
	}

	value := map[string]interface{}{
		FieldPassportStatus: status,
	}
	return orm.TableWithUpdate(ctx, exec, TablePassportName, value, map[string]interface{}{FieldPassportId: id})
}

func PassportWithUpdatePassword(ctx context.Context, exec mysql.Exec, salt, saltPassword string, id string) (int64, error) {
	if salt == "" || saltPassword == "" {
		return 0, fmt.Errorf("salt/salt_password is nil")
	}

	value := map[string]interface{}{
		FieldPassportSalt:         salt,
		FieldPassportSaltPassword: saltPassword,
	}
	return orm.TableWithUpdate(ctx, exec, TablePassportName, value, map[string]interface{}{FieldPassportId: id})
}

func PassportWithDeleteOne(ctx context.Context, exec mysql.Exec, id string) (int64, error) {
	if id == "" {
		return 0, fmt.Errorf("id is nil")
	}
	return orm.TableWithDelete(ctx, exec, TablePassportName, map[string]interface{}{FieldPassportId: id})
}

func PassportWithInsertOne(ctx context.Context, exec mysql.Exec, passport *Passport) (int64, error) {
	if passport == nil {
		return 0, fmt.Errorf("passport is nil")
	}

	value := map[string]interface{}{
		FieldPassportId:           passport.Id,
		FieldPassportAccount:      passport.Account,
		FieldPassportCode:         passport.Code,
		FieldPassportSalt:         passport.Salt,
		FieldPassportSaltPassword: passport.SaltPassword,
		FieldPassportEmail:        passport.Email,
		FieldPassportPhone:        passport.Phone,
		FieldPassportStatus:       passport.Status,
	}
	return orm.TableWithInsertOne(ctx, exec, TablePassportName, value)
}

type Passport struct {
	Id              string         `json:"id"`
	Account         string         `json:"account"`
	Code            sql.NullString `json:"code"`
	Salt            string         `json:"salt"`
	SaltPassword    string         `json:"salt_password"`
	Email           sql.NullString `json:"email"`
	Phone           sql.NullString `json:"phone"`
	Status          int32          `json:"status"`
	CreateTimestamp time.Time      `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime   `json:"modify_timestamp"`
}

const (
	TablePassportName = "passport"

	FieldPassportId              = "id"
	FieldPassportAccount         = "account"
	FieldPassportCode            = "code"
	FieldPassportSalt            = "salt"
	FieldPassportSaltPassword    = "salt_password"
	FieldPassportEmail           = "email"
	FieldPassportPhone           = "phone"
	FieldPassportStatus          = "status"
	FieldPassportCreateTimestamp = "create_timestamp"
	FieldPassportModifyTimestamp = "modify_timestamp"
)

var PassportFields = []string{
	FieldPassportId,
	FieldPassportAccount,
	FieldPassportCode,
	FieldPassportSalt,
	FieldPassportSaltPassword,
	FieldPassportEmail,
	FieldPassportPhone,
	FieldPassportStatus,
	FieldPassportCreateTimestamp,
	FieldPassportModifyTimestamp,
}

// type NullJson struct {
// 	Json  map[string]any
// 	Valid bool
// }

// func (nj NullJson) Value() (driver.Value, error) {
// 	if !nj.Valid {
// 		return nil, nil
// 	}

// 	if len(nj.Json) == 0 {
// 		return nil, nil
// 	}
// 	if nj.Json == nil {
// 		nj.Json = make(map[string]any, 8)
// 	}

// 	j, err := json.Marshal(nj.Json)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return driver.Value([]byte(j)), nil
// }

// func (nj *NullJson) Scan(src any) error {
// 	var source []byte
// 	switch src.(type) {
// 	case []uint8:
// 		source = []byte(src.([]uint8))
// 	case nil:
// 		return nil
// 	default:
// 		return errors.New("incompatible type for ValueInterfaceJSON")
// 	}

// 	nj.Json = make(map[string]interface{}, 8)
// 	err := json.Unmarshal(source, &nj.Json)
// 	if err != nil {
// 		return err
// 	}
// 	nj.Valid = true
// 	return nil
// }

// func (nj *NullJson) String() string {
// 	buf, _ := json.Marshal(nj)
// 	return string(buf)
// }
