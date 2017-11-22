package order

import (
	"github.com/less-leg/dbmodel"
	"github.com/less-leg/sql"
	"github.com/less-leg/test_model"
	"time"
)

type OrderDirection string

type relatedOrderEntity struct {
	Name    OrderField
	Version OrderField
}

type orderBuilder struct {
	fields []OrderField
}

func Desc(field OrderField) OrderDirection {
	return OrderDirection("")
}

type OrderField interface {
	sql.Field
}

type orderField struct {
	Column string
}

type OrderIdField interface {
	OrderField
	Equal(int) sql.Condition
}

type OrderOwnerField interface {
	OrderField

	Name() UserNameField
	Version() UserVersionField
}

type UserNameField interface {
	Like(v string, vs ...string) sql.Condition
}

type UserVersionField interface {
	NotIn(v int, vs ...int) sql.Condition
}

func Id() OrderIdField {
	return nil
}

type orderStartField struct {
	OrderField
}

func (t orderStartField) Equal(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Field: nil,
		Values:sql.MultiTimes(v),
	}
}

NotEqual(time.Time) sql.Condition
Greater(time.Time) sql.Condition
GreaterOrEqual(time.Time) sql.Condition
Lesser(time.Time) sql.Condition
LesserOrEqual(time.Time) sql.Condition
In(...time.Time) sql.Condition
NotIn(...time.Time) sql.Condition
IsNull() sql.Condition
IsNotNull() sql.Condition
Between(time.Time, time.Time) sql.Condition
NotBetween(time.Time, time.Time) sql.Condition

type OrderStartField interface {
	OrderField
	Equal(time.Time) sql.Condition
	NotEqual(time.Time) sql.Condition
	Greater(time.Time) sql.Condition
	GreaterOrEqual(time.Time) sql.Condition
	Lesser(time.Time) sql.Condition
	LesserOrEqual(time.Time) sql.Condition
	In(...time.Time) sql.Condition
	NotIn(...time.Time) sql.Condition
	IsNull() sql.Condition
	IsNotNull() sql.Condition
	Between(time.Time, time.Time) sql.Condition
	NotBetween(time.Time, time.Time) sql.Condition
}

type OrderTaxFreeField interface {
	OrderField
	Equal(dbmodel.Confirmation) sql.Condition
	NotEqual(dbmodel.Confirmation) sql.Condition
	In(...dbmodel.Confirmation) sql.Condition
	NotIn(...dbmodel.Confirmation) sql.Condition
	IsNull() sql.Condition
	IsNotNull() sql.Condition
}

func TaxFree() OrderTaxFreeField {
	return nil
}

var allFields = []OrderField{Id, Owner()}

func Owner() OrderOwnerField {
	return nil
}

func Order(fields ...OrderField) orderBuilder {
	if len(fields) == 0 {
		return orderBuilder{fields: allFields}
	} else {
		return orderBuilder{fields: fields}
	}
	return orderBuilder{}
}

func (b orderBuilder) Where(...sql.Condition) orderBuilder {
	return b
}

func (b orderBuilder) GroupBy(column ...OrderField) orderBuilder {
	return b
}

func (b orderBuilder) Having(...sql.Condition) orderBuilder {
	return b
}

func (b orderBuilder) OrderBy(...sql.SortOrder) orderBuilder {
	return b
}

func (b orderBuilder) Get() ([]test_model.Order, error) {
	return nil, nil
}

func (b orderBuilder) GetPtr() ([]*test_model.Order, error) {
	return nil, nil
}
