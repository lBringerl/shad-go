//go:build !solution

package externalsort

import (
	"fmt"
	"io"
)

const ReadBufferSize = 100

type LineReaderImpl struct {
	reader    io.Reader
	buffer    []byte
	accBuffer []byte
	EOFError  bool
}

func findEndOfLine(line []byte) int {
	for i, b := range line {
		if b == '\n' {
			return i
		}
	}
	return -1
}

func (r *LineReaderImpl) ReadLine() (string, error) {
	endOfLineIdx := findEndOfLine(r.accBuffer)
	for endOfLineIdx == -1 && !r.EOFError {
		n, err := r.reader.Read(r.buffer)
		if err != nil {
			switch err {
			case io.EOF:
				r.EOFError = true
			default:
				return "", fmt.Errorf("r.reader.Read: %w", err)
			}
		}

		buffEndOfLineIdx := findEndOfLine(r.buffer[:n])
		if buffEndOfLineIdx != -1 {
			endOfLineIdx = len(r.accBuffer) + buffEndOfLineIdx
		}
		r.accBuffer = append(r.accBuffer, r.buffer[:n]...)
	}

	if len(r.accBuffer) == 0 && r.EOFError {
		return "", io.EOF
	}

	if r.EOFError {
		switch endOfLineIdx {
		case len(r.accBuffer) - 1:
			res := string(r.accBuffer[:endOfLineIdx])
			r.accBuffer = make([]byte, 0)
			return res, nil
		case -1:
			res := string(r.accBuffer)
			r.accBuffer = make([]byte, 0)
			return res, io.EOF
		default:
		}
	}

	res := string(r.accBuffer[:endOfLineIdx])
	r.accBuffer = r.accBuffer[endOfLineIdx+1:]

	return res, nil
}

func NewReader(r io.Reader) LineReader {
	return &LineReaderImpl{
		reader:    r,
		buffer:    make([]byte, ReadBufferSize),
		accBuffer: make([]byte, 0),
		EOFError:  false,
	}
}

// type LineWriterImpl struct {
// 	writer io.Writer
// }

// func () Write(l string) error {}

func NewWriter(w io.Writer) LineWriter {
	panic("implement me")
}

func Merge(w LineWriter, readers ...LineReader) error {
	panic("implement me")
}

func Sort(w io.Writer, in ...string) error {
	panic("implement me")
}
