package logger

import (
	"fmt"
	"time"
)

func Log(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	message = fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
	if len(message) > 126 {
		message = message[:123] + "..."
	}
	fmt.Println(message)
}
