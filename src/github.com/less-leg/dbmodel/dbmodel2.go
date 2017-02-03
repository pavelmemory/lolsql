package dbmodel

import (
	"time"
	"encoding/json"
	"bytes"
)

type Handsome struct {
	Login       string `lolsql:"id[true] column[USER_LOGIN]"`
	Password    *int64 `lolsql:"column[SECRET]"`
	DateOfBirth *time.Time
	Salary      float32
}

func (*Handsome)TableName() string {
	return "USER"
}

type DjangoAdminLog struct {
	Id            int          `lolsql:"id[true]"`
	ActionTime    time.Time    `lolsql:"column[action_time]"`
	ObjectId      *string      `lolsql:"column[object_id]"`
	ObjectRepr    string       `lolsql:"column[object_repr]"`
	ActionFlag    int16        `lolsql:"column[action_flag]"`
	ChangeMessage string       `lolsql:"column[change_message]"`
	ContentTypeId *int         `lolsql:"column[content_type_id]"`
	UserId        int          `lolsql:"column[user_id]"`
	AskPassword   *Confirmation        `lolsql:"column[CONFIRMATION]"`
}

var buffer = &bytes.Buffer{}
var encoder = func() *json.Encoder { e := json.NewEncoder(buffer);e.SetIndent("  ", ""); return e }()

func (this *DjangoAdminLog) String() string {
	buffer.Reset()
	encoder.Encode(this)
	return buffer.String()
}

func (*DjangoAdminLog)TableName() string {
	return "django_admin_log"
}
