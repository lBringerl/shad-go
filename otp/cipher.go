//go:build !solution

package otp

import (
	"fmt"
	"io"
)

const buffSize = 10

type StreamSypherReader struct {
	prng    io.Reader
	r       io.Reader
	rBuf    []byte
	prngBuf []byte
}

func (r *StreamSypherReader) Read(p []byte) (n int, err error) {
	n = 0
	prngBuf := r.prngBuf
	rBuf := r.rBuf

	for {
		if n+buffSize > len(p) {
			rBuf = make([]byte, len(p)-n)
		}
		count, err := r.r.Read(rBuf)
		if count != len(prngBuf) {
			prngBuf = make([]byte, count)
		}
		_, _ = r.prng.Read(prngBuf)
		for i := range count {
			p[n] = rBuf[i] ^ prngBuf[i]
			n++
			if n >= len(p) {
				return n, nil
			}
		}
		if err != nil {
			switch err {
			case io.EOF:
				return n, io.EOF
			default:
				return n, fmt.Errorf("r.r.Read: %w", err)
			}
		}
	}
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &StreamSypherReader{
		prng:    prng,
		r:       r,
		rBuf:    make([]byte, buffSize),
		prngBuf: make([]byte, buffSize),
	}
}

type StreamSypherWriter struct {
	prng    io.Reader
	w       io.Writer
	wBuf    []byte
	prngBuf []byte
}

func (w *StreamSypherWriter) Write(p []byte) (n int, err error) {
	n = 0
	wBuf := w.wBuf

	for n < len(p) {
		_, _ = w.prng.Read(w.prngBuf)
		for i := range buffSize {
			wBuf[i] = p[n] ^ w.prngBuf[i]
			n++
			if n >= len(p) {
				tmp := make([]byte, n%buffSize)
				copy(tmp, wBuf[:n%buffSize])
				wBuf = tmp
				break
			}
		}
		count, err := w.w.Write(wBuf)
		if err != nil {
			return n - (buffSize - count), err
		}
	}
	return n, err
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &StreamSypherWriter{
		prng:    prng,
		w:       w,
		wBuf:    make([]byte, buffSize),
		prngBuf: make([]byte, buffSize),
	}
}
