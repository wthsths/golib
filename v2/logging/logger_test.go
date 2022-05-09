package gl_logging

import (
	"fmt"
	"runtime"
	"testing"
)

func Test_Logger(t *testing.T) {
	logBase := `/tmp`
	if runtime.GOOS == "windows" {
		logBase = `C:\logs`
	}

	Init("temp-svc", logBase)

	logger := NewLogger("some-title", "session-id")

	for i := 0; i < 10000; i++ {
		logger.Infof(fmt.Sprint(i))
	}
}
