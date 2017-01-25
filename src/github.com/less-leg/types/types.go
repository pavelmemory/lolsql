package types

import (
	"github.com/less-leg/utils"
	"strings"
	"strconv"
	"time"
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

func DefineAmountStrings(vnext []string) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountStringPtrs(values []*string) ConditionConstant {
	switch len(values) {
	case 0:
		return Null
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInts(vnext []int) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountIntPtrs(vnext []*int) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt8s(vnext []int8) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt8Ptrs(vnext []*int8) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt16s(vnext []int16) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt16Ptrs(vnext []*int16) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt32s(vnext []int32) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt32Ptrs(vnext []*int32) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt64s(vnext []int64) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountInt64Ptrs(vnext []*int64) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountFloat32s(vnext []float32) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountFloat32Ptrs(vnext []*float32) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountFloat64s(vnext []float64) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountFloat64Ptrs(vnext []*float64) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountTime_Times(vnext []time.Time) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func DefineAmountTime_TimePtrs(vnext []*time.Time) ConditionConstant {
	switch len(vnext) + 1 {
	case 0:
		panic("Unreachable code")
	case 1:
		return Single
	default:
		return Multi
	}
}

func renderEqualsNull(condition interface{}) string {
	return condition.(SelectColumn).Column() + " is null"
}

func renderEqualsNotNull(condition interface{}) string {
	return condition.(SelectColumn).Column() + " is not null"
}

///////////////////////////////////////////////////////////
// string
///////////////////////////////////////////////////////////
func renderEqualsString(condition interface{}) string {
	cstring := condition.(LolConditionString)
	return cstring.Column() + " = " + utils.Quote(cstring.Values()[0])
}

func renderEqualsStringPtr(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	return cstring.Column() + " = " + utils.Quote(*cstring.Values()[0])
}

func renderNotEqualsString(condition interface{}) string {
	cstring := condition.(LolConditionString)
	return cstring.Column() + " <> " + utils.Quote(cstring.Values()[0])
}

func renderNotEqualsStringPtr(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	return cstring.Column() + " <> " + utils.Quote(*cstring.Values()[0])
}

func renderLikeString(condition interface{}) string {
	cstring := condition.(LolConditionString)
	return cstring.Column() + " like " + utils.Quote(cstring.Values()[0])
}

func renderLikeStringPtr(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	return cstring.Column() + " like " + utils.Quote(*cstring.Values()[0])
}

func renderNotLikeString(condition interface{}) string {
	cstring := condition.(LolConditionString)
	return cstring.Column() + " not like " + utils.Quote(cstring.Values()[0])
}

func renderNotLikeStringPtr(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	return cstring.Column() + " not like " + utils.Quote(*cstring.Values()[0])
}

func renderEqualsStrings(condition interface{}) string {
	cstring := condition.(LolConditionString)
	return cstring.Column() + " in (" + strings.Join(utils.QuoteAll(cstring.Values()...), ", ") + ")"
}

func renderEqualsStringPtrs(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	return cstring.Column() + " in (" + strings.Join(utils.QuoteAllPtrs(cstring.Values()...), ", ") + ")"
}

func renderNotEqualsStrings(condition interface{}) string {
	cstring := condition.(LolConditionString)
	return cstring.Column() + " not in (" + strings.Join(utils.QuoteAll(cstring.Values()...), ", ") + ")"
}

func renderNotEqualsStringPtrs(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	return cstring.Column() + " not in (" + strings.Join(utils.QuoteAllPtrs(cstring.Values()...), ", ") + ")"
}

func renderLikeStrings(condition interface{}) string {
	cstring := condition.(LolConditionString)
	vstr := make([]string, 0, len(cstring.Values()))
	for _, vptr := range cstring.Values() {
		vstr = append(vstr, cstring.Column() + " like " + utils.Quote(vptr))
	}
	return strings.Join(vstr, " or ")
}

func renderLikeStringPtrs(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	vstr := make([]string, 0, len(cstring.Values()))
	for _, vptr := range cstring.Values() {
		if vptr != nil {
			vstr = append(vstr, cstring.Column() + " like " + utils.Quote(*vptr))
		}
	}
	return strings.Join(vstr, " or ")
}

func renderNotLikeStrings(condition interface{}) string {
	cstring := condition.(LolConditionString)
	vstr := make([]string, 0, len(cstring.Values()))
	for _, vptr := range cstring.Values() {
		vstr = append(vstr, cstring.Column() + " not like " + utils.Quote(vptr))
	}
	return strings.Join(vstr, " and ")
}

func renderNotLikeStringPtrs(condition interface{}) string {
	cstring := condition.(LolConditionStringPtr)
	vstr := make([]string, 0, len(cstring.Values()))
	for _, vptr := range cstring.Values() {
		if vptr != nil {
			vstr = append(vstr, cstring.Column() + " not like " + utils.Quote(*vptr))
		}
	}
	return strings.Join(vstr, " and ")
}

///////////////////////////////////////////////////////////
// int
///////////////////////////////////////////////////////////
func renderEqualsInt(condition interface{}) string {
	cond := condition.(LolConditionInt)
	return cond.Column() + " = " + utils.Quote(strconv.Itoa(cond.Values()[0]))
}

func renderEqualsIntPtr(condition interface{}) string {
	cond := condition.(LolConditionIntPtr)
	return cond.Column() + " = " + utils.Quote(strconv.Itoa(*cond.Values()[0]))
}

func renderNotEqualsInt(condition interface{}) string {
	cond := condition.(LolConditionInt)
	return cond.Column() + " <> " + utils.Quote(strconv.Itoa(cond.Values()[0]))
}

func renderNotEqualsIntPtr(condition interface{}) string {
	cond := condition.(LolConditionIntPtr)
	return cond.Column() + " <> " + utils.Quote(strconv.Itoa(*cond.Values()[0]))
}

func renderEqualsInts(condition interface{}) string {
	cond := condition.(LolConditionInt)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.IntToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsIntPtrs(condition interface{}) string {
	cond := condition.(LolConditionIntPtr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.IntPtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInts(condition interface{}) string {
	cond := condition.(LolConditionInt)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.IntToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsIntPtrs(condition interface{}) string {
	cond := condition.(LolConditionIntPtr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.IntPtrToString(cond.Values()...)...), ", ") + ")"
}
///////////////////////////////////////////////////////////
// int8
///////////////////////////////////////////////////////////
func renderEqualsInt8(condition interface{}) string {
	cond := condition.(LolConditionInt8)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderEqualsInt8Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt8Ptr)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderNotEqualsInt8(condition interface{}) string {
	cond := condition.(LolConditionInt8)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderNotEqualsInt8Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt8Ptr)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderEqualsInt8s(condition interface{}) string {
	cond := condition.(LolConditionInt8)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int8ToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsInt8Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt8Ptr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int8PtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt8s(condition interface{}) string {
	cond := condition.(LolConditionInt8)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int8ToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt8Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt8Ptr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int8PtrToString(cond.Values()...)...), ", ") + ")"
}

///////////////////////////////////////////////////////////
// int16
///////////////////////////////////////////////////////////
func renderEqualsInt16(condition interface{}) string {
	cond := condition.(LolConditionInt16)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderEqualsInt16Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt16Ptr)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderNotEqualsInt16(condition interface{}) string {
	cond := condition.(LolConditionInt16)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderNotEqualsInt16Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt16Ptr)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderEqualsInt16s(condition interface{}) string {
	cond := condition.(LolConditionInt16)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int16ToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsInt16Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt16Ptr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int16PtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt16s(condition interface{}) string {
	cond := condition.(LolConditionInt16)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int16ToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt16Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt16Ptr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int16PtrToString(cond.Values()...)...), ", ") + ")"
}

///////////////////////////////////////////////////////////
// int32
///////////////////////////////////////////////////////////
func renderEqualsInt32(condition interface{}) string {
	cond := condition.(LolConditionInt32)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderEqualsInt32Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt32Ptr)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderNotEqualsInt32(condition interface{}) string {
	cond := condition.(LolConditionInt32)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderNotEqualsInt32Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt32Ptr)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderEqualsInt32s(condition interface{}) string {
	cond := condition.(LolConditionInt32)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int32ToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsInt32Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt32Ptr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int32PtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt32s(condition interface{}) string {
	cond := condition.(LolConditionInt32)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int32ToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt32Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt32Ptr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int32PtrToString(cond.Values()...)...), ", ") + ")"
}

///////////////////////////////////////////////////////////
// int64
///////////////////////////////////////////////////////////
func renderEqualsInt64(condition interface{}) string {
	cond := condition.(LolConditionInt64)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderEqualsInt64Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt64Ptr)
	return cond.Column() + " = " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderNotEqualsInt64(condition interface{}) string {
	cond := condition.(LolConditionInt64)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(cond.Values()[0]), 10))
}

func renderNotEqualsInt64Ptr(condition interface{}) string {
	cond := condition.(LolConditionInt64Ptr)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatInt(int64(*cond.Values()[0]), 10))
}

func renderEqualsInt64s(condition interface{}) string {
	cond := condition.(LolConditionInt64)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int64ToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsInt64Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt64Ptr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Int64PtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt64s(condition interface{}) string {
	cond := condition.(LolConditionInt64)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int64ToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsInt64Ptrs(condition interface{}) string {
	cond := condition.(LolConditionInt64Ptr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Int64PtrToString(cond.Values()...)...), ", ") + ")"
}

///////////////////////////////////////////////////////////
// float32
///////////////////////////////////////////////////////////
func renderEqualsFloat32(condition interface{}) string {
	cond := condition.(LolConditionFloat32)
	return cond.Column() + " = " + utils.Quote(strconv.FormatFloat(float64(cond.Values()[0]), 'f', -1, 32))
}

func renderEqualsFloat32Ptr(condition interface{}) string {
	cond := condition.(LolConditionFloat32Ptr)
	return cond.Column() + " = " + utils.Quote(strconv.FormatFloat(float64(*cond.Values()[0]), 'f', -1, 32))
}

func renderNotEqualsFloat32(condition interface{}) string {
	cond := condition.(LolConditionFloat32)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatFloat(float64(cond.Values()[0]), 'f', -1, 32))
}

func renderNotEqualsFloat32Ptr(condition interface{}) string {
	cond := condition.(LolConditionFloat32Ptr)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatFloat(float64(*cond.Values()[0]), 'f', -1, 32))
}

func renderEqualsFloat32s(condition interface{}) string {
	cond := condition.(LolConditionFloat32)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Float32ToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsFloat32Ptrs(condition interface{}) string {
	cond := condition.(LolConditionFloat32Ptr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Float32PtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsFloat32s(condition interface{}) string {
	cond := condition.(LolConditionFloat32)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Float32ToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsFloat32Ptrs(condition interface{}) string {
	cond := condition.(LolConditionFloat32Ptr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Float32PtrToString(cond.Values()...)...), ", ") + ")"
}

///////////////////////////////////////////////////////////
// float64
///////////////////////////////////////////////////////////
func renderEqualsFloat64(condition interface{}) string {
	cond := condition.(LolConditionFloat64)
	return cond.Column() + " = " + utils.Quote(strconv.FormatFloat(float64(cond.Values()[0]), 'f', -1, 64))
}

func renderEqualsFloat64Ptr(condition interface{}) string {
	cond := condition.(LolConditionFloat64Ptr)
	return cond.Column() + " = " + utils.Quote(strconv.FormatFloat(float64(*cond.Values()[0]), 'f', -1, 64))
}

func renderNotEqualsFloat64(condition interface{}) string {
	cond := condition.(LolConditionFloat64)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatFloat(float64(cond.Values()[0]), 'f', -1, 64))
}

func renderNotEqualsFloat64Ptr(condition interface{}) string {
	cond := condition.(LolConditionFloat64Ptr)
	return cond.Column() + " <> " + utils.Quote(strconv.FormatFloat(float64(*cond.Values()[0]), 'f', -1, 64))
}

func renderEqualsFloat64s(condition interface{}) string {
	cond := condition.(LolConditionFloat64)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Float64ToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsFloat64Ptrs(condition interface{}) string {
	cond := condition.(LolConditionFloat64Ptr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.Float64PtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsFloat64s(condition interface{}) string {
	cond := condition.(LolConditionFloat64)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Float64ToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsFloat64Ptrs(condition interface{}) string {
	cond := condition.(LolConditionFloat64Ptr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.Float64PtrToString(cond.Values()...)...), ", ") + ")"
}

///////////////////////////////////////////////////////////
// time
///////////////////////////////////////////////////////////
func renderEqualsTime(condition interface{}) string {
	cond := condition.(LolConditionTime)
	return cond.Column() + " = " + utils.Quote(cond.Values()[0].Format("2006-01-02 15:04:05"))
}

func renderEqualsTimePtr(condition interface{}) string {
	cond := condition.(LolConditionTimePtr)
	return cond.Column() + " = " + utils.Quote((*cond.Values()[0]).Format("2006-01-02 15:04:05"))
}

func renderNotEqualsTime(condition interface{}) string {
	cond := condition.(LolConditionTime)
	return cond.Column() + " <> " + utils.Quote(cond.Values()[0].Format("2006-01-02 15:04:05"))
}

func renderNotEqualsTimePtr(condition interface{}) string {
	cond := condition.(LolConditionTimePtr)
	return cond.Column() + " <> " + utils.Quote((*cond.Values()[0]).Format("2006-01-02 15:04:05"))
}

func renderEqualsTimes(condition interface{}) string {
	cond := condition.(LolConditionTime)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.TimeToString(cond.Values()...)...), ", ") + ")"
}

func renderEqualsTimePtrs(condition interface{}) string {
	cond := condition.(LolConditionTimePtr)
	return cond.Column() + " in (" + strings.Join(utils.QuoteAll(utils.TimePtrToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsTimes(condition interface{}) string {
	cond := condition.(LolConditionTime)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.TimeToString(cond.Values()...)...), ", ") + ")"
}

func renderNotEqualsTimePtrs(condition interface{}) string {
	cond := condition.(LolConditionTimePtr)
	return cond.Column() + " not in (" + strings.Join(utils.QuoteAll(utils.TimePtrToString(cond.Values()...)...), ", ") + ")"
}

var ConditionRenderingMap = map[ConditionConstant]map[string]func(interface{}) string {
	Equals | Null: {
		"string":renderEqualsNull,

		"int":renderEqualsNull,
		"int8":renderEqualsNull,
		"int16":renderEqualsNull,
		"int32":renderEqualsNull,
		"int64":renderEqualsNull,

		"float32": renderEqualsNull,
		"float64": renderEqualsNull,

		"time.Time": renderEqualsNull,
		"*string":renderEqualsNull,

		"*int":renderEqualsNull,
		"*int8":renderEqualsNull,
		"*int16":renderEqualsNull,
		"*int32":renderEqualsNull,
		"*int64":renderEqualsNull,

		"*float32": renderEqualsNull,
		"*float64": renderEqualsNull,

		"*time.Time": renderEqualsNull,
	},

	Equals | Not | Null:{
		"string": renderEqualsNotNull,

		"int": renderEqualsNotNull,
		"int8": renderEqualsNotNull,
		"int16": renderEqualsNotNull,
		"int32": renderEqualsNotNull,
		"int64": renderEqualsNotNull,

		"float32": renderEqualsNotNull,
		"float64": renderEqualsNotNull,

		"time.Time": renderEqualsNotNull,

		"*string": renderEqualsNotNull,

		"*int": renderEqualsNotNull,
		"*int8": renderEqualsNotNull,
		"*int16": renderEqualsNotNull,
		"*int32": renderEqualsNotNull,
		"*int64": renderEqualsNotNull,

		"*float32": renderEqualsNotNull,
		"*float64": renderEqualsNotNull,

		"*time.Time": renderEqualsNotNull,
	},

	Equals | Single: {
		"string":renderEqualsString,

		"int":renderEqualsInt,
		"int8":renderEqualsInt8,
		"int16":renderEqualsInt16,
		"int32":renderEqualsInt32,
		"int64":renderEqualsInt64,

		"float32": renderEqualsFloat32,
		"float64": renderEqualsFloat64,

		"time.Time": renderEqualsTime,

		"*string":renderEqualsStringPtr,

		"*int":renderEqualsIntPtr,
		"*int8":renderEqualsInt8Ptr,
		"*int16":renderEqualsInt16Ptr,
		"*int32":renderEqualsInt32Ptr,
		"*int64":renderEqualsInt64Ptr,

		"*float32": renderEqualsFloat32Ptr,
		"*float64": renderEqualsFloat64Ptr,

		"*time.Time": renderEqualsTimePtr,
	},

	Not | Equals | Single: {
		"string":renderNotEqualsString,

		"int":renderNotEqualsInt,
		"int8":renderNotEqualsInt8,
		"int16":renderNotEqualsInt16,
		"int32":renderNotEqualsInt32,
		"int64":renderNotEqualsInt64,

		"float32": renderNotEqualsFloat32,
		"float64": renderNotEqualsFloat64,

		"time.Time": renderNotEqualsTime,

		"*string":renderNotEqualsStringPtr,

		"*int":renderNotEqualsIntPtr,
		"*int8":renderNotEqualsInt8Ptr,
		"*int16":renderNotEqualsInt16Ptr,
		"*int32":renderNotEqualsInt32Ptr,
		"*int64":renderNotEqualsInt64Ptr,

		"*float32": renderNotEqualsFloat32Ptr,
		"*float64": renderNotEqualsFloat64Ptr,

		"*time.Time": renderNotEqualsTimePtr,
	},

	Like | Single: {
		"string":renderLikeString,
		"*string":renderLikeStringPtr,
	},

	Not | Like | Single: {
		"string":renderNotLikeString,
		"*string":renderNotLikeStringPtr,
	},

	Equals | Multi: {
		"string":renderEqualsStrings,

		"int":renderEqualsInts,
		"int8":renderEqualsInt8s,
		"int16":renderEqualsInt16s,
		"int32":renderEqualsInt32s,
		"int64":renderEqualsInt64s,

		"float32": renderEqualsFloat32s,
		"float64": renderEqualsFloat64s,

		"time.Time": renderEqualsTimes,

		"*string":renderEqualsStringPtrs,

		"*int":renderEqualsIntPtrs,
		"*int8":renderEqualsInt8Ptrs,
		"*int16":renderEqualsInt16Ptrs,
		"*int32":renderEqualsInt32Ptrs,
		"*int64":renderEqualsInt64Ptrs,

		"*float32": renderEqualsFloat32Ptrs,
		"*float64": renderEqualsFloat64Ptrs,

		"*time.Time": renderEqualsTimePtrs,
	},

	Not | Equals | Multi: {
		"string":renderNotEqualsStrings,

		"int":renderNotEqualsInts,
		"int8":renderNotEqualsInt8s,
		"int16":renderNotEqualsInt16s,
		"int32":renderNotEqualsInt32s,
		"int64":renderNotEqualsInt64s,

		"float32": renderNotEqualsFloat32s,
		"float64": renderNotEqualsFloat64s,

		"time.Time": renderNotEqualsTimes,

		"*string":renderNotEqualsStringPtrs,

		"*int":renderNotEqualsIntPtrs,
		"*int8":renderNotEqualsInt8Ptrs,
		"*int16":renderNotEqualsInt16Ptrs,
		"*int32":renderNotEqualsInt32Ptrs,
		"*int64":renderNotEqualsInt64Ptrs,

		"*float32": renderNotEqualsFloat32Ptrs,
		"*float64": renderNotEqualsFloat64Ptrs,

		"*time.Time": renderNotEqualsTimePtrs,
	},

	Like | Multi: {
		"string":renderLikeStrings,
		"*string":renderLikeStringPtrs,
	},

	Not | Like | Multi: {
		"string":renderNotLikeStrings,
		"*string":renderNotLikeStringPtrs,
	},
}