package utils

import (
	"fmt"
	"log/syslog"
	"math"
	"math/rand"
	"strings"
	"time"
)

func RandInterval(min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
}

func Round(val float64, places int) float64 {
	f := math.Pow10(places)
	return float64(int64(val*f+0.5)) / f
}

func EncodePassword(password string) string {
	password = fmt.Sprintf("%s|%s", password, "")
	s := AesEncode(password)
	return s
}

func DecodePassword(password string) string {
	s := AesDecode(password)
	ary := strings.Split(s, "|")
	return ary[0]
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
