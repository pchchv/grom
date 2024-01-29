package grom

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func TestResponseWriterWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := ResponseWriter(&appResponseWriter{ResponseWriter: rec})

	assert.Equal(t, rw.Written(), false)

	n, err := rw.Write([]byte("Hello world"))
	assert.Equal(t, n, 11)
	assert.NoError(t, err)

	assert.Equal(t, n, 11)
	assert.Equal(t, rec.Code, rw.StatusCode())
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "Hello world")
	assert.Equal(t, rw.Size(), 11)
	assert.Equal(t, rw.Written(), true)
}

func TestResponseWriterWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := ResponseWriter(&appResponseWriter{ResponseWriter: rec})

	rw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, rec.Code, rw.StatusCode())
	assert.Equal(t, rec.Code, http.StatusNotFound)
}

func TestResponseWriterHijack(t *testing.T) {
	hijackable := &hijackableResponse{}
	rw := ResponseWriter(&appResponseWriter{ResponseWriter: hijackable})
	hijacker, ok := rw.(http.Hijacker)
	assert.True(t, ok)
	_, _, err := hijacker.Hijack()
	assert.NoError(t, err)
	assert.True(t, hijackable.Hijacked)
}

func TestResponseWriterHijackNotOK(t *testing.T) {
	rw := ResponseWriter(&appResponseWriter{ResponseWriter: httptest.NewRecorder()})
	_, _, err := rw.Hijack()
	assert.Error(t, err)
}
