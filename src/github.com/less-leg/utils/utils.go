package utils

import (
	"log"
	"strings"
	"strconv"
	"time"
)

var InvalidImportProtection = 0

func PanicIf(err error) {
	if err != nil {
		log.Panicln(err.Error())
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

func DoubleQuote(strs ...string) []string {
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

