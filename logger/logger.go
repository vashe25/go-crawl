package logger

import (
	"fmt"
	"time"
)

func Log(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	if len(message) > 126 {
		message = fmt.Sprintf("[%s]\n%s", time.Now().Format("2006-01-02 15:04:05"), message)
	} else {
		message = fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
	}
	fmt.Println(message)
}
