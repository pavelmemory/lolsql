package types

import (
	"github.com/less-leg/utils"
	"strings"
	"strconv"
	"time"
	"reflect"
)

type SelectColumn interface {
	Column() string
}

type LolCondition interface {
	SelectColumn
	Render() string
	And(LolCondition) LolCondition
	Or(LolCondition) LolCondition
	Next() LolCondition
	SetNext(LolCondition)
}

type LolConditionString interface {
	LolCondition
	Values() []string
}

type LolConditionStringPtr interface {
	LolCondition
	Values() []*string
}

type LolConditionInt interface {
	LolCondition
	Values() []int
}

type LolConditionIntPtr interface {
	LolCondition
	Values() []*int
}

type LolConditionInt8 interface {
	LolCondition
	Values() []int8
}

type LolConditionInt8Ptr interface {
	LolCondition
	Values() []*int8
}

type LolConditionInt16 interface {
	LolCondition
	Values() []int16
}

type LolConditionInt16Ptr interface {
	LolCondition
	Values() []*int16
}

type LolConditionInt32 interface {
	LolCondition
	Values() []int32
}

type LolConditionInt32Ptr interface {
	LolCondition
	Values() []*int32
}

type LolConditionInt64 interface {
	LolCondition
	Values() []int64
}

type LolConditionInt64Ptr interface {
	LolCondition
	Values() []*int64
}

type LolConditionFloat32 interface {
	LolCondition
	Values() []float32
}

type LolConditionFloat32Ptr interface {
	LolCondition
	Values() []*float32
}

type LolConditionFloat64 interface {
	LolCondition
	Values() []float64
}

type LolConditionFloat64Ptr interface {
	LolCondition
	Values() []*float64
}

type LolConditionTime interface {
	LolCondition
	Values() []time.Time
}

type LolConditionTimePtr interface {
	LolCondition
	Values() []*time.Time
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

func DefineAmount(v1 interface{}, vnext interface{}) ConditionConstant {
	value1 := reflect.ValueOf(v1)
	switch value1.Kind() {
	case reflect.Array, reflect.Slice:
		if !reflect.ValueOf(vnext).IsValid() {
			switch utils.Length(v1) {
			case 0:
				return Null
			case 1:
				return Single
			default:
				return Multi
			}
		} else {
			panic("Incorrect usage of the method. Expected: DefineAmount(val, vals) | DefineAmount(vals, nil)")
		}
	default:
		if utils.Length(vnext) + 1 == 1 {
			return Single
		} else {
			return Multi
		}
	}
}

func renderEqualsNull(condition interface{}) string {
	return condition.(SelectColumn).Column() + " is null"
}

func renderEqualsNotNull(condition interface{}) string {
	return condition.(SelectColumn).Column() + " is not null"
}

func renderEquals(condition interface{}) string {
	return renderSingle(condition, " = ")
}

func renderNotEquals(condition interface{}) string {
	return renderSingle(condition, " <> ")
}

func renderLike(condition interface{}) string {
	return renderSingle(condition, " like ")
}

func renderNotLike(condition interface{}) string {
	return renderSingle(condition, " not like ")
}

func renderMultiEquals(condition interface{}) string {
	return renderMulti(condition, " in ")
}

func renderMultiNotEquals(condition interface{}) string {
	return renderMulti(condition, " not in ")
}

func renderSingle(condition interface{}, comparator string) string {
	switch tcondition := condition.(type) {
	case LolConditionString:
		return tcondition.Column() + comparator + utils.Quote(tcondition.Values()[0])
	case LolConditionStringPtr:
		return tcondition.Column() + comparator + utils.Quote(*tcondition.Values()[0])

	case LolConditionInt:
		return tcondition.Column() + comparator + strconv.Itoa(tcondition.Values()[0])
	case LolConditionIntPtr:
		return tcondition.Column() + comparator + strconv.Itoa(*tcondition.Values()[0])
	case LolConditionInt8:
		return tcondition.Column() + comparator + strconv.FormatInt(int64(tcondition.Values()[0]), 10)
	case LolConditionInt8Ptr:
		return tcondition.Column() + comparator + strconv.FormatInt(int64(*tcondition.Values()[0]), 10)
	case LolConditionInt16:
		return tcondition.Column() + comparator + strconv.FormatInt(int64(tcondition.Values()[0]), 10)
	case LolConditionInt16Ptr:
		return tcondition.Column() + comparator + strconv.FormatInt(int64(*tcondition.Values()[0]), 10)
	case LolConditionInt32:
		return tcondition.Column() + comparator + strconv.FormatInt(int64(tcondition.Values()[0]), 10)
	case LolConditionInt32Ptr:
		return tcondition.Column() + comparator + strconv.FormatInt(int64(*tcondition.Values()[0]), 10)
	case LolConditionInt64:
		return tcondition.Column() + comparator + strconv.FormatInt(tcondition.Values()[0], 10)
	case LolConditionInt64Ptr:
		return tcondition.Column() + comparator + strconv.FormatInt(*tcondition.Values()[0], 10)

	case LolConditionFloat32:
		return tcondition.Column() + comparator + strconv.FormatFloat(float64(tcondition.Values()[0]), 'f', -1, 32)
	case LolConditionFloat32Ptr:
		return tcondition.Column() + comparator + strconv.FormatFloat(float64(*tcondition.Values()[0]), 'f', -1, 32)
	case LolConditionFloat64:
		return tcondition.Column() + comparator + strconv.FormatFloat(float64(tcondition.Values()[0]), 'f', -1, 64)
	case LolConditionFloat64Ptr:
		return tcondition.Column() + comparator + strconv.FormatFloat(float64(*tcondition.Values()[0]), 'f', -1, 64)

	case LolConditionTime:
		return tcondition.Column() + comparator + utils.Quote(tcondition.Values()[0].Format("2006-01-02 15:04:05"))
	case LolConditionTimePtr:
		return tcondition.Column() + comparator + utils.Quote((*tcondition.Values()[0]).Format("2006-01-02 15:04:05"))

	default:
		// TODO: Custom User types must provide mechanism to propagate them to sql query
		panic("Not supported type")
	}
}

func renderMulti(condition interface{}, comparator string) string {
	switch tcondition := condition.(type) {
	case LolConditionString:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.QuoteAll(tcondition.Values()...), ", ") + ")"
	case LolConditionStringPtr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.QuoteAllPtrs(tcondition.Values()...), ", ") + ")"

	case LolConditionInt:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.IntToString(tcondition.Values()...), ", ") + ")"
	case LolConditionIntPtr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.IntPtrToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt8:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int8ToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt8Ptr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int8PtrToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt16:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int16ToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt16Ptr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int16PtrToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt32:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int32ToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt32Ptr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int32PtrToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt64:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int64ToString(tcondition.Values()...), ", ") + ")"
	case LolConditionInt64Ptr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Int64PtrToString(tcondition.Values()...), ", ") + ")"

	case LolConditionFloat32:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Float32ToString(tcondition.Values()...), ", ") + ")"
	case LolConditionFloat32Ptr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Float32PtrToString(tcondition.Values()...), ", ") + ")"
	case LolConditionFloat64:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Float64ToString(tcondition.Values()...), ", ") + ")"
	case LolConditionFloat64Ptr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.Float64PtrToString(tcondition.Values()...), ", ") + ")"

	case LolConditionTime:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.QuoteAll(utils.TimeToString(tcondition.Values()...)...), ", ") + ")"
	case LolConditionTimePtr:
		return tcondition.Column() + comparator + "(" + strings.Join(utils.QuoteAll(utils.TimePtrToString(tcondition.Values()...)...), ", ") + ")"

	default:
		panic("Not supported type")
	}
}

func renderLikeStrings(condition interface{}) string {
	return renderMultiLike(condition, " like ", " or ")
}

func renderNotLikeStrings(condition interface{}) string {
	return renderMultiLike(condition, " not like ", " and ")
}

func renderMultiLike(condition interface{}, comparator, logic string) string {
	switch cstring := condition.(type) {
	case LolConditionString:
		vstr := make([]string, 0, len(cstring.Values()))
		for _, vptr := range cstring.Values() {
			vstr = append(vstr, cstring.Column() + comparator + utils.Quote(vptr))
		}
		return strings.Join(vstr, logic)
	case LolConditionStringPtr:
		vstr := make([]string, 0, len(cstring.Values()))
		for _, vptr := range cstring.Values() {
			if vptr != nil {
				vstr = append(vstr, cstring.Column() + comparator + utils.Quote(*vptr))
			}
		}
		return strings.Join(vstr, logic)
	}
	panic("Not supported type. Must be string or pointer to string")
}

var ConditionRenderingMap = map[ConditionConstant]func(interface{}) string{
	Equals | Null: renderEqualsNull,
	Equals | Not | Null: renderEqualsNotNull,
	Equals | Single: renderEquals,
	Not | Equals | Single: renderNotEquals,
	Like | Single: renderLike,
	Not | Like | Single: renderNotLike,

	Equals | Multi: renderMultiEquals,
	Not | Equals | Multi: renderMultiNotEquals,

	Like | Multi: renderLikeStrings,
	Not | Like | Multi: renderNotLikeStrings,
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