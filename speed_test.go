package grom

import "net/http"

// Null response writer
type NullWriter struct{}

func (w *NullWriter) Header() http.Header {
	return nil
}

func (w *NullWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func (w *NullWriter) WriteHeader(statusCode int) {}

// Types used by any/all frameworks:
type RouterBuilder func(namespaces []string, resources []string) http.Handler

// Benchmarks for gocraft/web:
type BenchContext struct {
	MyField string
}

type BenchContextB struct {
	*BenchContext
}

type BenchContextC struct {
	*BenchContextB
}
