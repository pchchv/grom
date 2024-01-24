package grom

import (
	"fmt"
	"runtime"
	"strings"
)

type nullPanicReporter struct{}

func (l nullPanicReporter) Panic(url string, err interface{}, stack string) {
	// no op
}

func init() {
	// This disables printing panics to stderr during testing, because that is very noisy,
	// and we purposefully test some panics.
	PanicHandler = nullPanicReporter{}
}

// callerInfo returns the caller's caller info.
func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}
