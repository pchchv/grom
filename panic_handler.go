package grom

// PanicReporter can receive panics that happen when serving
// a request and report them to a log of some sort.
type PanicReporter interface {
	// Panic is called with the URL of the request,
	// the result of calling recover, and the stack.
	Panic(url string, err interface{}, stack string)
}
