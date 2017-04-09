package utils

import (
	"log"
	"os"
	"reflect"
	"strings"
)

var InvalidImportProtection = 0

func PanicIfNotNil(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func RecreateDirectory(fileDir string) string {
	PanicIfNotNil(os.RemoveAll(fileDir))
	PanicIfNotNil(os.MkdirAll(fileDir, os.ModePerm))
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

func Suffix(suff string, strs ...string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		modif = append(modif, suff+str)
	}
	return modif
}

func Prefix(prefix string, strs ...string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		modif = append(modif, str+prefix)
	}
	return modif
}

func EncloseWith(suffix, prefix string, strs ...string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		modif = append(modif, suffix+str+prefix)
	}
	return modif
}

func EncloseWithPtrs(suffix, prefix string, strs ...*string) []string {
	modif := make([]string, 0, len(strs))
	for _, str := range strs {
		if str != nil {
			modif = append(modif, suffix+*str+prefix)
		}
	}
	return modif
}
