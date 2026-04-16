package config

import (
	"os"
	"strings"
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

func IsRunningInContainer() bool {
	if os.Getenv("EKHOES_RUNNING_IN_CONTAINER") != "" {
		return os.Getenv("EKHOES_RUNNING_IN_CONTAINER") == "true"
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil {
		c := string(data)
		if strings.Contains(c, "docker") ||
			strings.Contains(c, "kubepods") ||
			strings.Contains(c, "containerd") {
			return true
		}
	}

	return false
}
