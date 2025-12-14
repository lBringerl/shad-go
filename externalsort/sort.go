//go:build !solution

package externalsort

import (
	"container/heap"
	"fmt"
	"io"
	"os"
)

const (
	ReadBufferSize = 100000
	WriteBufferSize
)

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

type LineWriterImpl struct {
	writer io.Writer
	buffer []byte
}

func (w *LineWriterImpl) Write(l string) error {
	if len(l)+1 > len(w.buffer) {
		w.buffer = make([]byte, len(l)+1)
	}
	copy(w.buffer, l)
	w.buffer[len(l)] = '\n'

	_, err := w.writer.Write(w.buffer[:len(l)+1])
	if err != nil {
		return fmt.Errorf("w.writer.Write: %w", err)
	}

	return nil
}

func NewWriter(w io.Writer) LineWriter {
	return &LineWriterImpl{
		writer: w,
		buffer: make([]byte, WriteBufferSize),
	}
}

type StringHeap []string

func (h StringHeap) Len() int { return len(h) }

func (h StringHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h StringHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *StringHeap) Push(x any) {
	*h = append(*h, x.(string))
}

func (h *StringHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Merge(w LineWriter, readers ...LineReader) error {
	readersMap := make(map[LineReader]*string)

	for _, reader := range readers {
		readersMap[reader] = nil
	}

	stringsHeap := make(StringHeap, 0)

	for len(readersMap) != 0 {
		for reader, lastLine := range readersMap {
			if lastLine == nil || stringsHeap.Len() == 0 || *lastLine <= stringsHeap[0] {
				line, err := reader.ReadLine()
				if err != nil {
					switch err {
					case io.EOF:
						delete(readersMap, reader)
					default:
						return fmt.Errorf("reader.ReadLine: %w", err)
					}
				}
				if line != "" || err != io.EOF {
					heap.Push(&stringsHeap, line)
					readersMap[reader] = &line
				}
			}
		}

		if stringsHeap.Len() > 0 {
			line := heap.Pop(&stringsHeap)
			err := w.Write(line.(string))
			if err != nil {
				return fmt.Errorf("lineWriter.Write: %w", err)
			}
		}
	}

	if stringsHeap.Len() > 0 {
		line := heap.Pop(&stringsHeap)
		err := w.Write(line.(string))
		if err != nil {
			return fmt.Errorf("lineWriter.Write: %w", err)
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	for _, filename := range in {
		file, err := os.OpenFile(filename, os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("os.OpenFile: %w", err)
		}
		defer file.Close()

		fileContent := make(StringHeap, 0)
		fileContentCap := cap(fileContent)

		lineReader := NewReader(file)
	ReadLoop:
		for {
			if fileContentCap != cap(fileContent) {
				fileContentCap = cap(fileContent)
				heap.Init(&fileContent)
			}
			line, err := lineReader.ReadLine()
			if line != "" || err != io.EOF {
				heap.Push(&fileContent, line)
			}
			if err != nil {
				switch err {
				case io.EOF:
					break ReadLoop
				default:
					return fmt.Errorf("lineReader.ReadLine: %w", err)
				}
			}
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("file.Seek: %w", err)
		}
		lineWriter := NewWriter(file)

		for fileContent.Len() != 0 {
			line := heap.Pop(&fileContent)

			err = lineWriter.Write(line.(string))
			if fileContent.Len() == 0 && line.(string) != "" {
				st, err := file.Stat()
				if err != nil {
					panic(err)
				}
				file.Truncate(st.Size() - 1)
			}
			if err != nil {
				return fmt.Errorf("lineWriter.Write: %w", err)
			}
		}
	}

	readers := make([]LineReader, 0)
	for _, filename := range in {
		file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("os.OpenFile: %w", err)
		}
		readers = append(readers, NewReader(file))
	}

	err := Merge(NewWriter(w), readers...)
	if err != nil {
		return fmt.Errorf("Merge: %w", err)
	}

	return nil
}
