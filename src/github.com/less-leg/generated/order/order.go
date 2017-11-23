package order

import (
	"github.com/less-leg/dbmodel"
	"github.com/less-leg/generated/common"
	"github.com/less-leg/parser"
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
	uow    UnitOfWork
	fields []OrderField
}

func Desc(field OrderField) OrderDirection {
	return OrderDirection("")
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

type OrderField interface {
	sql.Field
}

type orderField struct{}

func (orderField) GetType() parser.TypeIdentity {
	return parser.TypeIdentity{Name: "Order", Package: "github.com/less-leg/test_model"}
}

type orderStartField struct {
	orderField
}

func Start() OrderStartField {
	return orderStartField{}
}

func (orderStartField) GetName() string {
	return "Start"
}

func (t orderStartField) Equal(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Equal,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) NotEqual(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotEqual,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) Greater(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Greater,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) GreaterOrEqual(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.GreaterOrEqual,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) Lesser(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Lesser,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) LesserOrEqual(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.LesserOrEqual,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) In(v ...time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.In,
		Values:              common.MultiTimeTime(v...),
	}
}

func (t orderStartField) NotIn(v ...time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotIn,
		Values:              common.MultiTimeTime(v...),
	}
}
func (t orderStartField) IsNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNull,
	}
}

func (t orderStartField) IsNotNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNotNull,
	}
}

func (t orderStartField) Like(vs ...time.Time) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Conjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.Like,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t orderStartField) LikeOr(vs ...time.Time) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Disjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.Like,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t orderStartField) NotLike(vs ...time.Time) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Conjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.NotLike,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t orderStartField) NotLikeOr(vs ...time.Time) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Disjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.NotLike,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t orderStartField) Between(v1 time.Time, v2 time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Between,
		Values:              []interface{}{v1, v2},
	}
}

func (t orderStartField) NotBetween(v1 time.Time, v2 time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotBetween,
		Values:              []interface{}{v1, v2},
	}
}

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

func Owner() OrderOwnerField {
	return nil
}

type UnitOfWork struct {
}

func (uow UnitOfWork) Order(fields ...OrderField) orderBuilder {
	return orderBuilder{uow: uow}
}

func Order(fields ...OrderField) orderBuilder {
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
