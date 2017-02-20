package common

import (
	"log"
)

func PanicOrLog(enablePanic bool, enableLog bool, err error) {
	if err != nil {
		if enablePanic {
			log.Panic(err)
		}
		if enableLog {
			log.Println(err)
		}
	}
}
