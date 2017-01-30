package generator

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/less-leg/utils"
	"github.com/less-leg/dbmodel"
	"fmt"
	"reflect"
	"github.com/less-leg/types"
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
	db       *sql.DB
	sql      string
	selector Selector
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

//func SetWithRetrieval(mapper func(interface{}, interface{}), retrieval func()) func(interface{}, interface{}) {
//	return func(to interface{}, from interface{}) {
//		mapper(to, retrieval())
//	}
//}

func (this *djangoAdminLogScanner) Select(field string) {
	this.selector.InitSelection(field)
}

type holder struct {
	tmp       interface{}
	propagate func(*dbmodel.DjangoAdminLog, interface{})
}

type Selector struct {
	holders    map[string]*holder
	fieldNames []string
}

func (this *Selector) InitSelection(fieldName string) {
	if this.holders == nil {
		this.holders = make(map[string]*holder)
	}
	if _, found := this.holders[fieldName]; found {
		panic("Incorrectly builded query string: Duplicate select field: " + fieldName)
	}

	var h *holder
	switch fieldName {
	case "Id":
		h = &holder{
			tmp:new(int),
			propagate:func(entity *dbmodel.DjangoAdminLog, tmp interface{}) {
				entity.Id = *(tmp.(*int))
			},
		}

	case "ActionTime":
		h = &holder{
			tmp:new(types.NullTime),
			propagate:func(entity *dbmodel.DjangoAdminLog, tmp interface{}) {
				t := tmp.(*types.NullTime)
				if t.Valid {
					entity.ActionTime = t.Time
				}
			},
		}

	case "ObjectId":
		h = &holder{tmp:new(sql.NullString), propagate:func(entity *dbmodel.DjangoAdminLog, tmp interface{}) {
			s := *(tmp.(*sql.NullString))
			if s.Valid {
				entity.ObjectId = &s.String
			}
		}}
	default:
		panic(fieldName)
	}

	this.holders[fieldName] = h
	this.fieldNames = append(this.fieldNames, fieldName)
}

func (this *Selector) Tmps() []interface{} {
	tmps := make([]interface{}, 0, len(this.fieldNames))
	for _, fieldName := range this.fieldNames {
		if holder, found := this.holders[fieldName]; found {
			tmps = append(tmps, holder.tmp)
		} else {
			panic("Selection builder has error: Holder for field was not found: " + fieldName)
		}
	}
	return tmps
}

func (this *Selector) Fetch() *dbmodel.DjangoAdminLog {
	entity := new(dbmodel.DjangoAdminLog)
	for _, fieldName := range this.fieldNames {
		if holder, found := this.holders[fieldName]; found {
			holder.propagate(entity, holder.tmp)
		} else {
			panic("Selection builder has error: Holder for field was not found: " + fieldName)
		}
	}
	return entity
}

func (this *djangoAdminLogScanner) Fucj() ([]*dbmodel.DjangoAdminLog, error) {
	rows, err := this.db.Query(this.sql)
	if err != nil {
		return nil, err
	}

	fetched := []*dbmodel.DjangoAdminLog{}
	for rows.Next() {
		if err = rows.Scan(this.selector.Tmps()...); err != nil {
			if errClose := rows.Close(); errClose != nil {
				return nil, fmt.Errorf("Cause: %s\nClose: %s", err.Error(), errClose.Error())
			}
			return nil, err
		}
		fetched = append(fetched, this.selector.Fetch())
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

func SelectDjangoAdminLog(db *sql.DB) ([]*dbmodel.DjangoAdminLog, error) {
	dalScaner := djangoAdminLogScanner{}
	dalScaner.db = db
	dalScaner.sql = "select Id, action_time, object_id from django_admin_log where id = 2 and action_flag = 1"
	dalScaner.Select("Id")
	dalScaner.Select("ActionTime")
	dalScaner.Select("ObjectId")
	return dalScaner.Fucj()

}