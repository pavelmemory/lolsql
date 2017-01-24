package utils

import (
	"log"
	"strings"
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