package logger

import (
	"log"
)

var isLogEnabled bool

func Config(enable bool) {
	isLogEnabled = enable
}

func Info(v string) {
	if isLogEnabled {
		log.Println("INFO: ", v)
	}
}

func Error(v string) {
	if isLogEnabled {
		log.Println("ERROR: ", v)
	}
}
