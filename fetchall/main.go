//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
)

const maxSize = 10000

func readContent(content io.ReadCloser, chunkSize int) ([]byte, error) {
	fullData := make([]byte, 0)
	buffer := make([]byte, chunkSize)

	for {
		n, err := content.Read(buffer)

		switch err {
		case nil:
		case io.EOF:
			return slices.Concat(fullData, buffer[:n]), nil
		default:
			return nil, fmt.Errorf("resp.Body.Read: %w", err)
		}

		fullData = slices.Concat(fullData, buffer[:n])
	}
}

type readFromUrlResponse struct {
	Url     string
	Data    string
	Err     error
	Elapsed int64
}

func readFromUrl(address string, ch chan readFromUrlResponse) {
	start := time.Now()

	resp, err := http.Get(address)
	if err != nil {
		ch <- readFromUrlResponse{
			Url:     address,
			Data:    "",
			Err:     fmt.Errorf("http.Get: %w", err),
			Elapsed: time.Since(start).Microseconds(),
		}
		return
	}

	data, err := readContent(resp.Body, maxSize)
	if err != nil {
		ch <- readFromUrlResponse{
			Url:     address,
			Data:    "",
			Err:     fmt.Errorf("readContent: %w", err),
			Elapsed: time.Since(start).Microseconds(),
		}
		return
	}

	ch <- readFromUrlResponse{
		Url:     address,
		Data:    string(data),
		Err:     nil,
		Elapsed: time.Since(start).Microseconds(),
	}
}

func main() {
	syncChannel := make(chan readFromUrlResponse)

	for _, url := range os.Args[1:] {
		go readFromUrl(url, syncChannel)
	}

	for range len(os.Args[1:]) {
		res := <-syncChannel
		fmt.Printf("%s: %v\n", res.Url, res.Elapsed)
	}
}
