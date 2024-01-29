package grom

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
)

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

func (h *hijackableResponse) Flush() {
	// no-op
}

func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func (h *hijackableResponse) CloseNotify() <-chan bool {
	return nil
}

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}
