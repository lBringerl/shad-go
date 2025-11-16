//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
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

func readFromUrl(address string) (string, error) {
	resp, err := http.Get(address)
	if err != nil {
		return "", fmt.Errorf("http.Get: %w", err)
	}

	data, err := readContent(resp.Body, maxSize)
	if err != nil {
		return "", fmt.Errorf("readContent: %w", err)
	}

	return string(data), nil
}

func main() {
	for _, url := range os.Args[1:] {
		data, err := readFromUrl(url)
		if err != nil {
			panic(err)
		}
		fmt.Printf("data: %v", data)
	}
}
