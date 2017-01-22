package person

import (
	"strings"
	"fmt"
	"strconv"
	. "github.com/less-leg/types"
)
type column interface {
	Column() string
}
//=====================SELECT=======================
type lol struct {
	selectColumns  []column
	whereInited bool
}
func Select(selects ...column) *lol {
	return &lol{selectColumns:selects}
}
func (this *lol) Render() string {
	if (len(this.selectColumns) == 0) {
		return "select id, name, address from books_author"
	}
	cols := make([]string, 0, len(this.selectColumns))
	for _, selectColumn := range this.selectColumns {
		cols = append(cols, selectColumn.Column())
	}
	return "select " + strings.Join(cols, ", ") + " from books_author"
}
func (this *lol) Where(cond LolCondition) *lolWhere {
	if this.whereInited {
		panic("Invalid usage of WHERE statement: double usage not supported.")
	}
	this.whereInited = true
	return &lolWhere{retrieval:this, condition:cond}
}
//=====================WHERE=======================
type lolWhere struct {
	retrieval *lol
	condition LolCondition
	next      []LolCondition
}
func (this *lolWhere) Render() string {
	if len(this.next) > 0 {
		conds := make([]string, 0, len(this.next))
		for _, cond := range this.next {
			conds = append(conds, cond.Render())
		}
		return this.retrieval.Render() + " where (" + this.condition.Render() + ") " + strings.Join(conds, " ")
	}
	return this.retrieval.Render() + " where " + this.condition.Render()
}
func (this *lolWhere) And(cond LolCondition) *lolWhere {
	if this.next == nil {
		this.next = make([]LolCondition, 0, 1)
	}
	this.next = append(this.next, &lolConditionAnd{condition: cond})
	return this
}
func (this *lolWhere) Or(cond LolCondition) *lolWhere {
	if this.next == nil {
		this.next = make([]LolCondition, 0, 1)
	}
	this.next = append(this.next, &lolConditionOr{condition: cond})
	return this
}
//=====================COLUMNS==========================

type idStub struct {
}

func Id() *idStub {
	return &idStub{}
}
func (_ *idStub) Column() string {
	return "id"
}

type nameStub struct { column }
var nameStubConst nameStub
func Name() *nameStub {return &nameStubConst}
func (*nameStub) Column() string {return "name"}

type firstNameStub struct{}
func FirstName() *firstNameStub {return &firstNameStub{}}
func (_ *firstNameStub) Column() string {return "first_name"}

type lastNameStub struct{}
func LastName() *lastNameStub {return &lastNameStub{}}
func (_ *lastNameStub) Column() string {return "last_name"}

type emailStub struct{}
func Email() *emailStub {return &emailStub{}}
func (_ *emailStub) Column() string {return "email"}

//==================================================
//=====================CONDITION====================
type idPerson struct {
	HasNext
	values    []*int
	checkNull bool
}
func (pId *idPerson) NotNull() LolCondition {
	if len(pId.values) > 0 {
		panic("Invalid usage: cannot check both equality and NULL check in one condition")
	}
	pId.checkNull = true
	return pId
}
func (this *idPerson) render() string {
	if this.values == nil || len(this.values) == 0 {
		if this.checkNull {
			return "id is not null"
		} else {
			return "id is null"
		}
	}
	if (len(this.values) == 1) {
		return fmt.Sprintf("id = %d", *(this.values[0]))
	}
	vstr := make([]string, 0, len(this.values))
	for _, vptr := range this.values {
		if vptr != nil {
			vstr = append(vstr, strconv.Itoa(*vptr))
		}
	}
	return fmt.Sprintf("id in (%s)", strings.Join(vstr, ", "))
}
func (this *idPerson) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}
func (this *idPerson) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}
func (this *idPerson) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}
func IdIs(value ...*int) LolCondition {
	return &idPerson{values:value}
}//=====================CONDITION====================
type salaryPerson struct {
	HasNext
	values    []*float64
	checkNull bool
}
func (this *salaryPerson) NotNull() LolCondition {
	if len(this.values) > 0 {
		panic("Invalid usage: cannot check both equality and NULL check in one condition")
	}
	this.checkNull = true
	return this
}
func (this *salaryPerson) render() string {
	if this.values == nil || len(this.values) == 0 {
		if this.checkNull {
			return "salary is not null"
		} else {
			return "salary is null"
		}
	}
	if (len(this.values) == 1) {
		return fmt.Sprintf("salary = %f", *(this.values[0]))
	}
	vstr := make([]string, 0, len(this.values))
	for _, vptr := range this.values {
		if vptr != nil {
			vstr = append(vstr, strconv.FormatFloat(*vptr, 'f', -1, 64))
		}
	}
	return fmt.Sprintf("salary in (%s)", strings.Join(vstr, ", "))
}
func (this *salaryPerson) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}
func (this *salaryPerson) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}
func (this *salaryPerson) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}
func SalaryIs(value ...*float64) LolCondition {
	return &salaryPerson{values:value}
}
//===============Name=========================================
type namePerson struct {
	HasNext
	values   []*string
	checkNot bool
}

func (this *namePerson) render() string {
	if this.values == nil || len(this.values) == 0 {
		if this.checkNot {
			return "name is not null"
		} else {
			return "name is null"
		}
	}
	if len(this.values) == 1 {
		if this.checkNot {
			return fmt.Sprintf("name <> '%s'", *(this.values[0]))
		} else {
			return fmt.Sprintf("name = '%s'", *(this.values[0]))
		}
	}
	vstr := make([]string, 0, len(this.values))
	for _, vptr := range this.values {
		if vptr != nil {
			vstr = append(vstr, "'" + *vptr + "'")
		} else {
			panic("NULL can't be used as one of the arguments in the list")
		}
	}
	if this.checkNot {
		return fmt.Sprintf("name not in (%s)", strings.Join(vstr, ", "))
	} else {
		return fmt.Sprintf("name in (%s)", strings.Join(vstr, ", "))
	}
}
func (this *namePerson) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}
func (this *namePerson) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}
func (this *namePerson) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}
func NameIs(values ...*string) *namePerson {
	return &namePerson{values:values}
}
func NameIsNot(values ...*string) *namePerson {
	return &namePerson{values:values, checkNot: true}
}
//==================================================
//=====================AND=======================
type lolConditionAnd struct {
	HasNext
	condition LolCondition
}
func (this *lolConditionAnd) render() string {
	return "and (" + this.condition.Render() + ")"
}
func (this *lolConditionAnd) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}
func (this *lolConditionAnd) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}
func (this *lolConditionAnd) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}
//===============================================
//=====================OR========================
type lolConditionOr struct {
	HasNext
	condition LolCondition
}
func (this *lolConditionOr) render() string {
	return "or (" + this.condition.Render() + ")"
}
func (this *lolConditionOr) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}
func (this *lolConditionOr) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}
func (this *lolConditionOr) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}
//===============================================



