//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func isASCII(b byte) bool {
	return b < 0x80
}

func isContinuationByte(b byte) bool {
	return (b & 0xC0) == 0x80
}

func isTwoBytesStart(b byte) bool {
	return (b & 0xE0) == 0xC0
}

func isThreeBytesStart(b byte) bool {
	return (b & 0xF0) == 0xE0
}

func isFourBytesStart(b byte) bool {
	return (b & 0xF8) == 0xF0
}

func handleMultibyte(resString *strings.Builder, continuationCounter int, expectedSeqLength int, input string, i int) {
	if continuationCounter < expectedSeqLength-1 {
		for range continuationCounter + 1 {
			resString.WriteRune(utf8.RuneError)
		}
	} else if continuationCounter == expectedSeqLength-1 {
		for j := range expectedSeqLength {
			resString.WriteByte(input[i+j])
		}
	} else {
		for range continuationCounter - expectedSeqLength {
			resString.WriteRune(utf8.RuneError)
		}
		for j := range expectedSeqLength {
			resString.WriteByte(input[i+j])
		}
	}
}

func Reverse(input string) string {
	var resString strings.Builder
	resString.Grow(len(input))

	continuationCounter := 0
	for i := len(input) - 1; i >= 0; i-- {
		continuationCounter = 0
		for ; i >= 0 && isContinuationByte(input[i]); i-- {
			continuationCounter++
		}

		if isTwoBytesStart(input[i]) {
			handleMultibyte(&resString, continuationCounter, 2, input, i)
		} else if isThreeBytesStart(input[i]) {
			handleMultibyte(&resString, continuationCounter, 3, input, i)
		} else if isFourBytesStart(input[i]) {
			handleMultibyte(&resString, continuationCounter, 4, input, i)
		} else if isASCII(input[i]) {
			resString.WriteByte(input[i])
			for range continuationCounter {
				resString.WriteRune(utf8.RuneError)
			}
		} else {
			resString.WriteRune(utf8.RuneError)
			for range continuationCounter {
				resString.WriteRune(utf8.RuneError)
			}
		}
	}

	return resString.String()
}
