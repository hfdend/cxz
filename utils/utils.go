package utils

import (
	"log/syslog"
	"math/rand"
	"strings"
	"time"
)

func RandInterval(min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
}

func ParseSyslogPriority(s string) syslog.Priority {
	var priority syslog.Priority
	switch strings.ToUpper(s) {
	default:
		priority = syslog.LOG_LOCAL0
	case "LOCAL0":
		priority = syslog.LOG_LOCAL0
	case "LOCAL1":
		priority = syslog.LOG_LOCAL1
	case "LOCAL2":
		priority = syslog.LOG_LOCAL2
	case "LOCAL3":
		priority = syslog.LOG_LOCAL3
	case "LOCAL4":
		priority = syslog.LOG_LOCAL4
	case "LOCAL5":
		priority = syslog.LOG_LOCAL5
	case "LOCAL6":
		priority = syslog.LOG_LOCAL6
	case "LOCAL7":
		priority = syslog.LOG_LOCAL7
	}
	return priority
}
