package tools

import (
	"fmt"
	"time"
)

func MydateType(logleve int, message string) string {

	var dateType string

	currentTime := time.Now()

	if logleve == 1 {
		dateType = fmt.Sprintf("%s [INFO] %s", currentTime.Format("2006.01.02 15:04:05"), message)
	} else if logleve == 2 {
		dateType = fmt.Sprintf("%s [ERROR] %s", currentTime.Format("2006.01.02 15:04:05"), message)
	}
	return dateType
}
