package gl_logging

import (
	"fmt"
	"testing"
)

func Test_Logger(t *testing.T) {
	Init("temp-svc", `C:\logs`)

	logger := NewLogger("some-title", "session-id")

	for i := 0; i < 10000; i++ {
		logger.Infof(fmt.Sprint(i))
	}
}
