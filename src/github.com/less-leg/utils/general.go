package utils

import "log"

func PanicIf(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}
