package main

import (
	"log"
)

func ErrLog(v ...interface{}) {
	TypeMonitorChan <- TypeErrNum
	log.Println(v)
}

func ExceptionLog(v ...interface{}) {
	TypeMonitorChan <- TypeExceptionNum
	log.Println(v)
}

func WarningLog(v ...interface{}) {
	TypeMonitorChan <- TypeWarningNum
	log.Println(v)
}