package types

import "reflect"

type LolDB struct{}

func GetTableName(entity interface{}) string {
	entityType := reflect.TypeOf(entity)
	if tableNameMethod, ok := entityType.MethodByName("TableName"); ok {
		return tableNameMethod.Func.Call(nil)[0].String()
	} else {
		return entityType.Name()
	}
}

type SelectColumn interface {
	Column() string
}

type LolCondition interface {
	Render() string
	And(LolCondition) LolCondition
	Or(LolCondition) LolCondition
	Next() LolCondition
	SetNext(LolCondition)
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