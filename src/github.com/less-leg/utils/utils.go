package utils

import (
	"log"
	"strings"
)

func PanicIf(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func DotToUnderscore(str string) string {
	return strings.Replace(str, ".", "_", -1)
}