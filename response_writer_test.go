package grom

import "net/http"

type hijackableResponse struct {
	Hijacked bool
}

func (h *hijackableResponse) Header() http.Header {
	return nil
}
func (h *hijackableResponse) Write(buf []byte) (int, error) {
	return 0, nil
}
func (h *hijackableResponse) WriteHeader(code int) {
	// no-op
}
