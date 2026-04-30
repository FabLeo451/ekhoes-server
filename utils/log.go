package utils

import (
	"ekhoes-server/config"
	"fmt"
	"log"
)

const (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	cyan   = "\033[36m" // azzurro
)

func Log(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.Printf("%s[INFO]%s %s", cyan, reset, msg)
}

func Debug(format string, a ...any) {
	if !config.Debug() {
		return
	}

	msg := fmt.Sprintf(format, a...)
	log.Printf("%s[DEBUG]%s %s", gray, reset, msg)
}

func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.Printf("%s[ERROR]%s %s", red, reset, msg)
}

func Err(err error) {
	log.Printf("%s[ERROR]%s %s", red, reset, err.Error())
}
