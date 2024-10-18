package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eviltomorrow/king/lib/codes"
	"github.com/eviltomorrow/king/lib/db/mysql"
	"github.com/eviltomorrow/king/lib/sqlutil"
	jsoniter "github.com/json-iterator/go"
)

func SchedulerRecordWithInsertOne(ctx context.Context, exec mysql.Exec, record *SchedulerRecord) (int64, error) {
	if record == nil {
		return 0, fmt.Errorf("scheduler record is nil")
	}

	value := map[string]interface{}{
		FieldSchedulerRecordId:          record.Id,
		FieldSchedulerRecordAlias:       record.Alias,
		FieldSchedulerRecordName:        record.Name,
		FieldSchedulerRecordDate:        record.Date,
		FieldSchedulerRecordServiceName: record.ServiceName,
		FieldSchedulerRecordFuncName:    record.FuncName,
		FieldSchedulerRecordStatus:      record.Status,
		FieldSchedulerRecordCode:        record.Code,
		FieldSchedulerRecordErrorMsg:    record.ErrorMsg,
	}
	return sqlutil.TableWithInsertOne(ctx, exec, TableSchedulerRecordName, value)
}

func SchedulerRecordWithUpdateStatus(ctx context.Context, exec mysql.Exec, status, code, errorMsg string, id string) (int64, error) {
	if status == "" {
		return 0, fmt.Errorf("status is invalid")
	}
	if code != codes.SUCCESS && errorMsg == "" {
		return 0, fmt.Errorf("code is invalid, msg: %s", errorMsg)
	}

	value := map[string]interface{}{
		FieldSchedulerRecordStatus: status,
		FieldSchedulerRecordCode:   code,
	}
	if errorMsg != "" {
		value[FieldSchedulerRecordErrorMsg] = errorMsg
	}
	return sqlutil.TableWithUpdate(ctx, exec, TableSchedulerRecordName, value, map[string]interface{}{FieldSchedulerRecordId: id})
}

func SchedulerRecordWithSelectOneByDateName(ctx context.Context, exec mysql.Exec, name, date string) (*SchedulerRecord, error) {
	record := SchedulerRecord{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&record.Id,
			&record.Alias,
			&record.Name,
			&record.Date,
			&record.ServiceName,
			&record.FuncName,
			&record.Status,
			&record.Code,
			&record.ErrorMsg,
			&record.ParentId,
			&record.CreateTimestamp,
			&record.ModifyTimestamp,
		)
	}
	if err := sqlutil.TableWithSelectOne(ctx, exec, TableSchedulerRecordName, schedulerRecordFields, map[string]interface{}{FieldSchedulerRecordName: name, fmt.Sprintf("DATE_FORMAT(`%s`, '%%Y-%%m-%%d')", FieldSchedulerRecordDate): date}, scan); err != nil {
		return nil, err
	}
	return &record, nil
}

func SchedulerRecordWithSelectOneById(ctx context.Context, exec mysql.Exec, id string) (*SchedulerRecord, error) {
	record := SchedulerRecord{}
	scan := func(row *sql.Row) error {
		return row.Scan(
			&record.Id,
			&record.Alias,
			&record.Name,
			&record.Date,
			&record.ServiceName,
			&record.FuncName,
			&record.Status,
			&record.Code,
			&record.ErrorMsg,
			&record.ParentId,
			&record.CreateTimestamp,
			&record.ModifyTimestamp,
		)
	}
	if err := sqlutil.TableWithSelectOne(ctx, exec, TableSchedulerRecordName, schedulerRecordFields, map[string]interface{}{FieldSchedulerRecordId: id}, scan); err != nil {
		return nil, err
	}
	return &record, nil
}

// Scheduler Record
type SchedulerRecord struct {
	Id              string         `json:"id"`
	Alias           string         `json:"alias"`
	Name            string         `json:"name"`
	Date            time.Time      `json:"date"`
	ServiceName     string         `json:"service_name"`
	FuncName        string         `json:"func_name"`
	Status          string         `json:"status"`
	Code            sql.NullString `json:"code"`
	ErrorMsg        sql.NullString `json:"error_msg"`
	ParentId        sql.NullString `json:"parent_id"`
	CreateTimestamp time.Time      `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime   `json:"modify_timestamp"`
}

func (s *SchedulerRecord) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s)
	return string(buf)
}

const (
	TableSchedulerRecordName = "scheduler_record"

	FieldSchedulerRecordId              = "id"
	FieldSchedulerRecordAlias           = "alias"
	FieldSchedulerRecordName            = "name"
	FieldSchedulerRecordDate            = "date"
	FieldSchedulerRecordServiceName     = "service_name"
	FieldSchedulerRecordFuncName        = "func_name"
	FieldSchedulerRecordStatus          = "status"
	FieldSchedulerRecordCode            = "code"
	FieldSchedulerRecordErrorMsg        = "error_msg"
	FieldSchedulerRecordParentId        = "parent_id"
	FieldSchedulerRecordCreateTimestamp = "create_timestamp"
	FieldSchedulerRecordModifyTimestamp = "modify_timestamp"
)

var schedulerRecordFields = []string{
	FieldSchedulerRecordId,
	FieldSchedulerRecordAlias,
	FieldSchedulerRecordName,
	FieldSchedulerRecordDate,
	FieldSchedulerRecordServiceName,
	FieldSchedulerRecordFuncName,
	FieldSchedulerRecordStatus,
	FieldSchedulerRecordCode,
	FieldSchedulerRecordErrorMsg,
	FieldSchedulerRecordParentId,
	FieldSchedulerRecordCreateTimestamp,
	FieldSchedulerRecordModifyTimestamp,
}
