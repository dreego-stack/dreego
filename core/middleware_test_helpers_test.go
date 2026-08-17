package core

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
)

type hijackableRecorder struct {
	*httptest.ResponseRecorder
}

func newHijackableRecorder() *hijackableRecorder {
	return &hijackableRecorder{ResponseRecorder: httptest.NewRecorder()}
}

func (hijackableRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

type pusherRecorder struct {
	*httptest.ResponseRecorder
}

func newPusherRecorder() *pusherRecorder {
	return &pusherRecorder{ResponseRecorder: httptest.NewRecorder()}
}

func (pusherRecorder) Push(string, *http.PushOptions) error {
	return nil
}

type readerFromRecorder struct {
	*httptest.ResponseRecorder
}

func newReaderFromRecorder() *readerFromRecorder {
	return &readerFromRecorder{ResponseRecorder: httptest.NewRecorder()}
}

func (r readerFromRecorder) ReadFrom(src io.Reader) (int64, error) {
	b, err := io.ReadAll(src)
	if err != nil {
		return 0, err
	}
	n, werr := r.ResponseRecorder.Write(b)
	return int64(n), werr
}
