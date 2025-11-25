//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func writeValue(builder *strings.Builder, val interface{}) {
	switch v := val.(type) {
	case int:
		builder.WriteString(strconv.Itoa(v))
	default:
		builder.WriteString(fmt.Sprint(v))
	}
}

func Sprintf(format string, args ...interface{}) string {
	var result strings.Builder
	result.Grow(len(format))

	var (
		iStart     int
		iEnd       = -1
		argCounter int
	)

	for i := 0; i < len(format); i++ {
		if format[i] == '{' {
			argCounter++

			iStart = i
			for ; i < len(format) && format[i] != '}'; i++ {

			}
			if i >= len(format) {
				panic("closing bracket not found")
			}

			result.WriteString(format[iEnd+1 : iStart])
			iEnd = i

			if i == iStart+1 {
				data := args[argCounter-1]
				writeValue(&result, data)
			} else {
				argIndexStr, err := strconv.Atoi(format[iStart+1 : i])
				if err != nil {
					panic(fmt.Sprintf("unexpected int conversion error %s", err.Error()))
				}

				data := args[argIndexStr]
				writeValue(&result, data)
			}
		}
	}
	result.WriteString(format[iEnd+1:])

	return result.String()
}
