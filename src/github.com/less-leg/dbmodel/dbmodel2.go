package dbmodel

import "time"

type Handsome struct {
	Login       string `lolsql:"id[true] column[USER_LOGIN]"`
	Password    *int64 `lolsql:"column[SECRET]"`
	DateOfBirth *time.Time
	Salary      float32
}

func (*Handsome)TableName() string {
	return "USER"
}