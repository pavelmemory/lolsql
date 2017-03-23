package types

import (
	"github.com/less-leg/utils"
	"strings"
)

type Columner interface {
	Column() string
}

type Fielder interface {
	Field() string
}

type FielderColumner interface {
	Fielder
	Columner
}

type LolCondition interface {
	Columner
	Render() string
	And(LolCondition) LolCondition
	Or(LolCondition) LolCondition
	Next() LolCondition
	SetNext(LolCondition)
	Parameters() []interface{}
}

func NewLolCondition(columner Columner, values []interface{}, operation ConditionConstant) LolCondition {
	return &LolConditionBase{
		Columner: columner,
		values: values,
		operation: operation,
	}
}

type LolConditionBase struct {
	Columner
	HasNext
	values   []interface{}
	operation ConditionConstant
}

func (this *LolConditionBase) Parameters() []interface{} {
	if this.Next() != nil {
		return append(this.values, this.Next().Parameters()...)
	}
	return this.values
}

func (this *LolConditionBase) render() string {
	if conditionSign, found := ConditionSignMap[this.operation]; found {
		return conditionSign(this.Column(), this.values)
	}
	panic("Not supported operation for: int")
}


func (this *LolConditionBase) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}


func (this *LolConditionBase) And(cond LolCondition) LolCondition {
	this.SetNext(&LolConditionAnd{LolCondition:cond})
	return this
}

func (this *LolConditionBase) Or(cond LolCondition) LolCondition {
	this.SetNext(&LolConditionAnd{LolCondition:cond})
	return this
}

type LolConditionAnd struct {
	HasNext
	LolCondition
}

func (this *LolConditionAnd) render() string {
	return "and (" + this.LolCondition.Render() + ")"
}

func (this *LolConditionAnd) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *LolConditionAnd) And(cond LolCondition) LolCondition {
	this.SetNext(&LolConditionAnd{LolCondition:cond})
	return this
}

func (this *LolConditionAnd) Or(cond LolCondition) LolCondition {
	this.SetNext(&LolConditionOr{LolCondition:cond})
	return this
}

func (this *LolConditionAnd) Next() LolCondition {
	return this.HasNext.Next()
}

func (this *LolConditionAnd) SetNext(n LolCondition) {
	this.HasNext.SetNext(n)
}

func (this *LolConditionAnd) Parameters() []interface{} {
	if this.Next() != nil {
		return append(this.LolCondition.Parameters(), this.Next().Parameters()...)
	} else {
		return this.LolCondition.Parameters()
	}
}

type LolConditionOr struct {
	HasNext
	LolCondition
}

func (this *LolConditionOr) render() string {
	return "or (" + this.LolCondition.Render() + ")"
}

func (this *LolConditionOr) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *LolConditionOr) And(cond LolCondition) LolCondition {
	this.SetNext(&LolConditionAnd{LolCondition:cond})
	return this
}

func (this *LolConditionOr) Or(cond LolCondition) LolCondition {
	this.SetNext(&LolConditionOr{LolCondition:cond})
	return this
}

func (this *LolConditionOr) Next() LolCondition {
	return this.HasNext.Next()
}

func (this *LolConditionOr) SetNext(n LolCondition) {
	this.HasNext.SetNext(n)
}

func (this *LolConditionOr) Parameters()[]interface{} {
	if this.Next() != nil {
		return append(this.LolCondition.Parameters(), this.Next().Parameters()...)
	} else {
		return this.LolCondition.Parameters()
	}
}

type HasNext struct {
	next LolCondition
}

func (this *HasNext) Next() LolCondition {
	return this.next
}

func (this *HasNext) SetNext(cond LolCondition) {
	prev := this.next
	next := this.next
	for next != nil {
		prev = next
		next = next.Next()
	}
	if prev != nil {
		prev.SetNext(cond)
	} else {
		this.next = cond
	}
}

type ConditionConstant int

const (
	None = ConditionConstant(0)
	_ = ConditionConstant(1 << (2 * iota))
	Null
	Single
	Multi
	Not
	Like
	Equals
)

func DefineConditionsAmount(count int) ConditionConstant {
	switch count {
	case 0:  return Null
	case 1:  return Single
	default: return Multi
	}
}

var ConditionSignMap = map[ConditionConstant]func(column string, values interface{}) string {
	Equals | Null: func(column string, values interface{}) string { return column + " is null" },
	Equals | Not | Null: func(column string, values interface{}) string { return column + " is not null" },

	Equals | Single: func(column string, values interface{}) string { return column + " = ?" },
	Not | Equals | Single: func(column string, values interface{}) string { return column + " <> ?" },

	Like | Single: func(column string, values interface{}) string { return column + " like ?" },
	Not | Like | Single: func(column string, values interface{}) string { return column + " not like ?" },

	Equals | Multi: func(column string, values interface{}) string { return column + " in(" + strings.Repeat("?,", utils.Length(values) - 1) + "?)"},
	Not | Equals | Multi: func(column string, values interface{}) string { return column + " not in(" + strings.Repeat("?,", utils.Length(values) - 1) + "?)"},

	Like | Multi: func(column string, values interface{}) string { return strings.Repeat(column + " like ? or ", utils.Length(values) - 1) + column + " like ?"},
	Not | Like | Multi: func(column string, values interface{}) string { return strings.Repeat(column + " not like ? and ", utils.Length(values) - 1) + column + " not like ?"},
}