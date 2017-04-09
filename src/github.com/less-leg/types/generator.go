package types

import (
	"github.com/less-leg/utils"
	"strings"
	"fmt"
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
	return &RootLolCondition{
		Columner:  columner,
		values:    values,
		operation: operation,
	}
}

type RootLolCondition struct {
	Columner
	HasNext
	values    []interface{}
	operation ConditionConstant
}

func (this *RootLolCondition) Parameters() []interface{} {
	if this.Next() != nil {
		return append(this.values, this.Next().Parameters()...)
	}
	return this.values
}

func (this *RootLolCondition) render() string {
	if conditionSign, found := ConditionSignMap[this.operation]; found {
		return conditionSign(this.Column(), this.values)
	}
	panic(fmt.Sprintf("Not supported operation for: %d", this.operation))
}

func (this *RootLolCondition) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *RootLolCondition) And(cond LolCondition) LolCondition {
	this.SetNext(NewAndCondition(cond))
	return this
}

func (this *RootLolCondition) Or(cond LolCondition) LolCondition {
	this.SetNext(NewOrCondition(cond))
	return this
}

func NewAndCondition(cond LolCondition) LolCondition {
	return &LolConditionLogical{LolCondition: cond, logicalCondition: "and"}
}

func NewOrCondition(cond LolCondition) LolCondition {
	return &LolConditionLogical{LolCondition: cond, logicalCondition: "or"}
}

type LolConditionLogical struct {
	HasNext
	LolCondition
	logicalCondition string
}

func (this *LolConditionLogical) Render() string {
	if this.Next() != nil {
		return this.logicalCondition + " (" + this.LolCondition.Render() + ") " + this.Next().Render()
	}
	return this.logicalCondition + " (" + this.LolCondition.Render() + ")"
}

func (this *LolConditionLogical) And(cond LolCondition) LolCondition {
	this.SetNext(NewAndCondition(cond))
	return this
}

func (this *LolConditionLogical) Or(cond LolCondition) LolCondition {
	this.SetNext(NewOrCondition(cond))
	return this
}

func (this *LolConditionLogical) Next() LolCondition {
	return this.HasNext.Next()
}

func (this *LolConditionLogical) SetNext(n LolCondition) {
	this.HasNext.SetNext(n)
}

func (this *LolConditionLogical) Parameters() []interface{} {
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
	// TODO: Should be extended with sql operators: BETWEEN, <, >, >=, <=
	None = ConditionConstant(0)
	_    = ConditionConstant(1 << (2 * iota))
	Null
	Single
	Multi
	Not
	Like
	Equals
)

func DefineConditionsAmount(count int) ConditionConstant {
	switch count {
	case 0:
		return Null
	case 1:
		return Single
	default:
		return Multi
	}
}

var ConditionSignMap = map[ConditionConstant]func(column string, values interface{}) string{
	// TODO: This map should be extended with sql operators: BETWEEN, <, >, >=, <=
	Equals | Null:       func(column string, values interface{}) string { return column + " is null" },
	Equals | Not | Null: func(column string, values interface{}) string { return column + " is not null" },

	Equals | Single:       func(column string, values interface{}) string { return column + " = ?" },
	Not | Equals | Single: func(column string, values interface{}) string { return column + " <> ?" },

	Like | Single:       func(column string, values interface{}) string { return column + " like ?" },
	Not | Like | Single: func(column string, values interface{}) string { return column + " not like ?" },

	Equals | Multi: func(column string, values interface{}) string {
		return column + " in(" + strings.Repeat("?,", utils.Length(values)-1) + "?)"
	},
	Not | Equals | Multi: func(column string, values interface{}) string {
		return column + " not in(" + strings.Repeat("?,", utils.Length(values)-1) + "?)"
	},

	Like | Multi: func(column string, values interface{}) string {
		return strings.Repeat(column+" like ? or ", utils.Length(values)-1) + column + " like ?"
	},
	Not | Like | Multi: func(column string, values interface{}) string {
		return strings.Repeat(column+" not like ? and ", utils.Length(values)-1) + column + " not like ?"
	},
}
