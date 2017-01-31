package utils

import (
	"log"
	"strings"
	"strconv"
	"time"
	"reflect"
	"os"
	"fmt"
)

var InvalidImportProtection = 0

func PanicIf(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func RecreateDirectory(fileDir string) string {
	err := os.RemoveAll(fileDir)
	if os.IsExist(err) {
		panic(fmt.Sprintf("Directory %s cannot be removed: %s", fileDir, err.Error()))
	}
	err = os.Mkdir(fileDir, os.ModePerm)
	if os.IsExist(err) {
		panic(fmt.Sprintf("Directory %s cannot be created: %s", fileDir, err.Error()))
	}
	return fileDir
}

func Length(elememts interface{}) int {
	value := reflect.ValueOf(elememts)
	if value.IsNil() {
		return 0
	} else {
		switch value.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			return value.Len()
		default:
			return 0
		}
	}
}

func DotToUnderscore(str string) string {
	return strings.Replace(str, ".", "_", -1)
}

func QuoteAll(strs ...string) []string {
	return EncloseWith("'", "'", strs...)
}

func QuoteAllPtrs(strs ...*string) []string {
	return EncloseWithPtrs("'", "'", strs...)
}

func Quote(strs string) string {
	return "'" + strs + "'"
}

func DoubleQuotes(strs ...string) []string {
	return EncloseWith("\"", "\"", strs...)
}

func Suffix(suff string, strs...string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		modif = append(modif, suff + str)
	}
	return modif
}

func Prefix(prefix string, strs...string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		modif = append(modif, str + prefix)
	}
	return modif
}

func EncloseWith(suffix, prefix string, strs ... string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		modif = append(modif, suffix + str + prefix)
	}
	return modif
}

func EncloseWithPtrs(suffix, prefix string, strs ... *string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		if str != nil {
			modif = append(modif, suffix + *str + prefix)
		}
	}
	return modif
}

func ToStrings(vals interface{}) [] string {
	if value := reflect.ValueOf(vals); value.IsValid() {
		var res []string
		switch value.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				res = append(res, ToString(value.Index(i).Interface()))
			}
		}
		return res
	}
	panic("Cannot be used as parameter")
}

var timeReflectType = reflect.TypeOf(time.Time{})

func ToString(vals interface{}) string {
	if value := reflect.ValueOf(vals); value.IsValid() {
		switch {
		case value.Kind() == reflect.Ptr:
			return valueToString(value.Elem())
		case reflect.TypeOf(vals) == timeReflectType:
			return vals.(time.Time).Format("2006-01-02 15:04:05")
		default:
			return valueToString(value)
		}
	}
	panic("Cannot be used as parameter")
}

func valueToString(value reflect.Value) string {
	if value.Type() == timeReflectType {
		return value.Interface().(time.Time).Format("2006-01-02 15:04:05")
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.String:
		return value.String()
	}
	panic("Not supported type for string creation: " + value.Kind().String())
}

func IntToString(ints ... int)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.Itoa(i))
	}
	return istrs
}

func IntPtrToString(ints ... *int)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.Itoa(*i))
		}
	}
	return istrs
}

func Int8ToString(ints ... int8)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.FormatInt(int64(i), 10))
	}
	return istrs
}

func Int8PtrToString(ints ... *int8)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.FormatInt(int64(*i), 10))
		}
	}
	return istrs
}

func Int16ToString(ints ... int16)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.FormatInt(int64(i), 10))
	}
	return istrs
}

func Int16PtrToString(ints ... *int16)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.FormatInt(int64(*i), 10))
		}
	}
	return istrs
}

func Int32ToString(ints ... int32)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.FormatInt(int64(i), 10))
	}
	return istrs
}

func Int32PtrToString(ints ... *int32)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.FormatInt(int64(*i), 10))
		}
	}
	return istrs
}

func Int64ToString(ints ... int64)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.FormatInt(i, 10))
	}
	return istrs
}

func Int64PtrToString(ints ... *int64)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.FormatInt(*i, 10))
		}
	}
	return istrs
}

func Float32ToString(ints ... float32)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.FormatFloat(float64(i), 'f', -1, 32))
	}
	return istrs
}

func Float32PtrToString(ints ... *float32)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.FormatFloat(float64(*i), 'f', -1, 32))
		}
	}
	return istrs
}

func Float64ToString(ints ... float64)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, strconv.FormatFloat(i, 'f', -1, 64))
	}
	return istrs
}

func Float64PtrToString(ints ... *float64)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, strconv.FormatFloat(*i, 'f', -1, 64))
		}
	}
	return istrs
}

func TimeToString(ints ... time.Time)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		istrs = append(istrs, i.Format("2006-01-02 15:04:05"))
	}
	return istrs
}

func TimePtrToString(ints ... *time.Time)[] string {
	istrs := make([]string, 0, len(ints))
	for _, i := range ints {
		if i != nil {
			istrs = append(istrs, i.Format("2006-01-02 15:04:05"))
		}
	}
	return istrs
}

func PrependString(v0 string, vnext []string) []string {
	buf := make([]string, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependInt(v0 int, vnext []int) []int {
	buf := make([]int, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependInt8(v0 int8, vnext []int8) []int8 {
	buf := make([]int8, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependInt16(v0 int16, vnext []int16) []int16 {
	buf := make([]int16, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependInt32(v0 int32, vnext []int32) []int32 {
	buf := make([]int32, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependInt64(v0 int64, vnext []int64) []int64 {
	buf := make([]int64, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependFloat32(v0 float32, vnext []float32) []float32 {
	buf := make([]float32, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependFloat64(v0 float64, vnext []float64) []float64 {
	buf := make([]float64, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
}

func PrependTime_Time(v0 time.Time, vnext []time.Time) []time.Time {
	buf := make([]time.Time, 0, 1 + len(vnext))
	buf = append(buf, v0)
	return append(buf, vnext...)
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