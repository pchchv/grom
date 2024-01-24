package grom

import "log"

// PanicReporter can receive panics that happen when serving
// a request and report them to a log of some sort.
type PanicReporter interface {
	// Panic is called with the URL of the request,
	// the result of calling recover, and the stack.
	Panic(url string, err interface{}, stack string)
}

type logPanicReporter struct {
	log *log.Logger
}

func (l logPanicReporter) Panic(url string, err interface{}, stack string) {
	l.log.Printf("PANIC\nURL: %v\nERROR: %v\nSTACK:\n%s\n", url, err, stack)
}
