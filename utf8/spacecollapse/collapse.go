//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
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

func handleMultibyte(resString *strings.Builder, expectedSeqLength int, input string, i int) int {
	continuationCounter := 0
	for j := 0; i+1+j < len(input) && isContinuationByte(input[i+1+j]); j++ {
		continuationCounter++
	}

	if continuationCounter == expectedSeqLength-1 {
		for j := i; j < i+expectedSeqLength; j++ {
			resString.WriteByte(input[j])
		}
	} else {
		for range continuationCounter + 1 {
			resString.WriteRune(utf8.RuneError)
		}
	}

	return i + continuationCounter
}

func CollapseSpaces(input string) string {
	var resString strings.Builder
	resString.Grow(len(input))

	for i := 0; i < len(input); i++ {
		j := 0
		for ; i+j < len(input) && (unicode.IsSpace(rune(input[i+j]))); j++ {
		}
		if j != 0 {
			resString.WriteByte(' ')
			i += j
			if i >= len(input) {
				break
			}
		}

		if isASCII(input[i]) {
			resString.WriteByte(input[i])
		} else if isTwoBytesStart(input[i]) {
			i = handleMultibyte(&resString, 2, input, i)
		} else if isThreeBytesStart(input[i]) {
			i = handleMultibyte(&resString, 3, input, i)
		} else if isFourBytesStart(input[i]) {
			i = handleMultibyte(&resString, 4, input, i)
		} else {
			resString.WriteRune(utf8.RuneError)
		}
	}

	return resString.String()
}
