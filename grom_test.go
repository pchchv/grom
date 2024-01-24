package grom

type nullPanicReporter struct{}

func (l nullPanicReporter) Panic(url string, err interface{}, stack string) {
	// no op
}
