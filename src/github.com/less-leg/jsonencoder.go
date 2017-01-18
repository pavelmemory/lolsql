package main

import (
	"reflect"
	"strconv"
	"fmt"
	"errors"
	"bytes"
)

var Constants = struct {
	Null             []byte
	Quote            []byte
	Comma            []byte
	Colon            []byte
	OpenCurlyBrace   []byte
	CloseCurlyBrace  []byte
	OpenSquareBrace  []byte
	CloseSquareBrace []byte
}{
	Null: []byte{'n', 'u', 'l', 'l'},
	Quote: []byte{'"'},
	Comma: []byte{','},
	Colon: []byte{':'},
	OpenCurlyBrace: []byte{'{'},
	CloseCurlyBrace: []byte{'}'},
	OpenSquareBrace : []byte{'['},
	CloseSquareBrace : []byte{']'},
}

func mainJson() {
	v := make(map[string]interface{})
	v["a"] = 10
	v["b"] = 101
	v["c"] = "qweqwe"
	v["d"] = nil
	v["e"] = []float32{3, 4, 9.0, -0.222}
	v["m"] = map[int]string{1: "one", 2: "two"}
	v["s"] = struct {
		Name string
		Age  int
	}{
		Name: "Pavel",
		Age:26,
	}
	v["o"] = map[bool]string{true: "True", false: "False"}
	fmt.Println(StringRef(v))

	fmt.Println(StringRef(struct {
		Name string
		Age  int
	}{
		Name: "Pavel",
		Age:26,
	}))

	fmt.Println(StringRef(10102))
}

func StringRef(v interface{}) string {
	return string(HighLevelValueAsString(reflect.ValueOf(v), false))
}

func HighLevelValueAsString(v reflect.Value, isKey bool) []byte {
	bStr := new(bytes.Buffer)

	switch v.Kind() {
	case reflect.Struct:
		bStr.Write(StructAsString(v))
	case reflect.Map:
		bStr.Write(MapAsString(v))
	case reflect.Interface, reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr:
		bStr.Write(InterfaceOrPointerAsString(v))
	case reflect.Slice, reflect.Array:
		bStr.Write(SliceOrArrayAsString(v))
	case
		reflect.Bool,
		reflect.String,
		reflect.Complex64, reflect.Complex128,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if isKey && (reflect.String != v.Kind() && reflect.Bool != v.Kind()) {
			bStr.Write(Constants.Quote)
			bStr.WriteString(LowLevelValueAsString(v))
			bStr.Write(Constants.Quote)
		} else {
			bStr.WriteString(LowLevelValueAsString(v))
		}

	case reflect.Chan, reflect.Func:
	// do not print any information about them
	default:
		panic(errors.New("Unhandled type for: " + v.String()))
	}

	return bStr.Bytes()
}

func MapAsString(v reflect.Value) []byte {
	if v.IsNil() {
		return Constants.Null
	} else {
		bStr := new(bytes.Buffer)
		bStr.Write(Constants.OpenCurlyBrace)
		mKeys := v.MapKeys()
		singleElement := true
		for _, mKey := range mKeys {
			if !singleElement {
				bStr.Write(Constants.Comma)
			} else {
				singleElement = false
			}

			bStr.Write(HighLevelValueAsString(mKey, true))
			bStr.Write(Constants.Colon)
			bStr.Write(HighLevelValueAsString(v.MapIndex(mKey), false))
		}
		bStr.Write(Constants.CloseCurlyBrace)
		return bStr.Bytes()
	}
}

func StructAsString(v reflect.Value) []byte {
	bStr := new(bytes.Buffer)
	bStr.Write(Constants.OpenCurlyBrace)
	singleElement := true
	for i := 0; i < v.NumField(); i += 1 {
		if !singleElement {
			bStr.Write(Constants.Comma)
		} else {
			singleElement = false
		}

		fieldType := v.Type().Field(i)
		bStr.Write(Constants.Quote)
		bStr.WriteString(fieldType.Name)
		bStr.Write(Constants.Quote)

		bStr.Write(Constants.Colon)
		bStr.Write(HighLevelValueAsString(v.Field(i), false))
	}
	bStr.Write(Constants.CloseCurlyBrace)
	return bStr.Bytes()
}

func InterfaceOrPointerAsString(v reflect.Value) []byte {
	if v.IsNil() {
		return Constants.Null
	} else {
		return HighLevelValueAsString(v.Elem(), false)
	}
}

func SliceOrArrayAsString(v reflect.Value) []byte {
	if v.IsNil() {
		return Constants.Null
	} else {
		bStr := new(bytes.Buffer)
		vLen := v.Len()
		singleElement := true
		bStr.Write(Constants.OpenSquareBrace)
		for i := 0; i < vLen; i += 1 {
			if !singleElement {
				bStr.Write(Constants.Comma)
			} else {
				singleElement = false
			}
			bStr.Write(HighLevelValueAsString(v.Index(i), false))
		}
		bStr.Write(Constants.CloseSquareBrace)
		return bStr.Bytes()
	}
}

// Convert any type into string
func LowLevelValueAsString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return string(strconv.AppendQuote(make([]byte, 0), v.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprint(v.Complex())
	case reflect.Bool:
		return string(strconv.AppendQuote(make([]byte, 0), strconv.FormatBool(v.Bool())))
	default:
		panic(errors.New("Unhandled type for: " + v.String()))
	}
}
