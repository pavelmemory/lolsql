package generator

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/less-leg/utils"
	"github.com/less-leg/dbmodel"
	"fmt"
	"github.com/less-leg/types"
	"reflect"
)

func Sdwef() {
	db, err := sql.Open("mysql", "root:root@/akeos")
	utils.PanicIf(err)
	defer db.Close()
	result, err := SelectDjangoAdminLog(db)
	utils.PanicIf(err)
	for _, v := range result {
		fmt.Println(v)
	}
}

type djangoAdminLogScanner struct {
	db        *sql.DB
	sql       string
	selectors []interface{}
	mappers []func(to interface{}, from interface{})
}

func Set(to interface{}, value interface{}) {
	valueOfTo := reflect.ValueOf(to)
	for valueOfTo.Kind() == reflect.Ptr {
		if valueOfTo.IsNil() {
			return
		}
		valueOfTo = valueOfTo.Elem()
	}
	valueOfValue := reflect.ValueOf(value)
	for valueOfValue.Kind() == reflect.Ptr {
		if valueOfValue.IsNil() {
			valueOfTo.Set(reflect.Zero(valueOfTo.Type()))
			return
		}
		valueOfValue = valueOfValue.Elem()
	}
	valueOfTo.Set(valueOfValue)
}

func SetWithRetrieval(mapper func(interface{}, interface{}), retrieval func()) func(interface{}, interface{}) {
	return func(to interface{}, from interface{}) {
		mapper(to, retrieval())
	}
}

func (this *djangoAdminLogScanner) Select(field interface{}, mapper func(interface{}, interface{})) {
	this.selectors = append(this.selectors, field)
	this.mappers = append(this.mappers, mapper)
}

func (this *djangoAdminLogScanner) Fucj(db *sql.DB) ([]*dbmodel.DjangoAdminLog, error) {
	dal := dbmodel.DjangoAdminLog{ObjectId:new(string), ContentTypeId:new(int)}
	objectId := sql.NullString{}
	actionTime := types.NullTime{}
	contentTypeId := sql.NullInt64{}

	this.Select(&dal.Id, Set)
	this.Select(&actionTime, SetWithRetrieval(Set, func() {
		if actionTime.Valid {
			dal.ActionTime = actionTime.Time
		}
	}))
	this.Select(&objectId, SetWithRetrieval(Set, func(){
		if objectId.Valid {
			dal.ObjectId = &objectId.String
		}
	}))
	this.Select(&dal.ObjectRepr, Set)
	this.Select(&dal.ActionFlag, Set)
	this.Select(&dal.ChangeMessage, Set)
	this.Select(&contentTypeId, SetWithRetrieval(Set, func(){
		if contentTypeId.Valid {
			contentTypeIdInt := int(contentTypeId.Int64)
			dal.ContentTypeId = &contentTypeIdInt
		}
	}))
	this.Select(&dal.UserId, Set)


	rows, err := this.db.Query(this.sql)
	if err != nil {
		return nil, err
	}

	fetched := []*dbmodel.DjangoAdminLog{}
	for rows.Next() {
		rows.Scan(this.selectors...)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return fetched, nil
}

func newDjangoAdminLog(fields []interface{}) *dbmodel.DjangoAdminLog {
	dal := &dbmodel.DjangoAdminLog{}
	dalValue := reflect.ValueOf(dal)
	for i := 0; i < dalValue.NumField(); i++ {
		Set(dalValue.Field(i).Interface(), fields[i])
	}
	return dal
}

// "select Id, action_time, object_id, object_repr, action_flag, change_message, content_type_id, user_id from django_admin_log where id = 2 and action_flag = 1"
// []interface{}{&dal.Id, &actionTime, &objectId, &dal.ObjectRepr, &dal.ActionFlag, &dal.ChangeMessage, &contentTypeId, &dal.UserId}
func (this *djangoAdminLogScanner) Scan() (fetched []*dbmodel.DjangoAdminLog, err error) {
	rows, err := this.db.Query(this.sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(this.selectors...)
		fetched = append(fetched, newDjangoAdminLog(this.selectors))


	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return fetched, nil
}

func SelectDjangoAdminLog(db *sql.DB) ([]*dbmodel.DjangoAdminLog, error) {
	rows, err := db.Query("select Id, action_time, object_id, object_repr, action_flag, change_message, content_type_id, user_id from django_admin_log where id = 2 and action_flag = 1")
	if err != nil {
		return nil, err
	}
	// TODO: get list of columns to propagate them into an entity object
	//rows.Columns()

	results := []*dbmodel.DjangoAdminLog{}
	for rows.Next() {
		dal := &dbmodel.DjangoAdminLog{}
		objectId := sql.NullString{}
		actionTime := types.NullTime{}
		contentTypeId := sql.NullInt64{}

		fields := []interface{}{&dal.Id, &actionTime, &objectId, &dal.ObjectRepr, &dal.ActionFlag, &dal.ChangeMessage, &contentTypeId, &dal.UserId}

		rows.Scan(fields...)

		if actionTime.Valid {
			fmt.Println("YES")
			dal.ActionTime = actionTime.Time
		} else {
			fmt.Println("NO")
		}
		if objectId.Valid {
			dal.ObjectId = &objectId.String
		}
		if contentTypeId.Valid {
			contentTypeIdInt := int(contentTypeId.Int64)
			dal.ContentTypeId = &contentTypeIdInt
		}

		results = append(results, dal)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return results, nil
}