package utils

import (
	"ekhoes-server/config"
	"fmt"
	"log"
)

func Log(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.Printf("[INFO] %s", msg)
}

func Debug(format string, a ...any) {
	if !config.Debug() {
		return
	}

	msg := fmt.Sprintf(format, a...)
	log.Printf("[DEBUG] %s", msg)
}

func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.Printf("[ERROR] %s", msg)
}

func Err(err error) {
	log.Printf("[ERROR] %s", err.Error())
}
