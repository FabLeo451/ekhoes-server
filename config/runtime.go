package config

import (
	"time"
)

type RuntimeStruct struct {
	StartTime time.Time
	Local     bool
	Database  string
	Cache     string
}

var Runtime RuntimeStruct

func Local() bool {
	return Runtime.Local
}
