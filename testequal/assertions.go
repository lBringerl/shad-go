//go:build !solution

package testequal

import "reflect"

func assertEqual(expected, actual interface{}) bool {
	switch expected.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64, string:
		expectedType := reflect.TypeOf(expected)
		actualType := reflect.TypeOf(actual)
		if expectedType != actualType || expected != actual {
			return false
		}
		return true
	case []int:
		expectedSlice := expected.([]int)
		actualSlice, ok := actual.([]int)
		expectedIsNil := expectedSlice == nil
		actualIsNil := actualSlice == nil
		if !ok || len(actualSlice) != len(expectedSlice) || actualIsNil != expectedIsNil {
			return false
		}
		for i := range actualSlice {
			if actualSlice[i] != expectedSlice[i] {
				return false
			}
		}
		return true
	case map[string]string:
		expectedMap := expected.(map[string]string)
		actualMap, ok := actual.(map[string]string)
		expectedIsNil := expectedMap == nil
		actualIsNil := actualMap == nil
		if !ok || len(actualMap) != len(expectedMap) || actualIsNil != expectedIsNil {
			return false
		}
		for key := range actualMap {
			if actualMap[key] != expectedMap[key] {
				return false
			}
		}
		return true
	case []byte:
		expectedBytes := expected.([]byte)
		actualBytes, ok := actual.([]byte)
		expectedIsNil := expectedBytes == nil
		actualIsNil := actualBytes == nil
		if !ok || len(actualBytes) != len(expectedBytes) || actualIsNil != expectedIsNil {
			return false
		}
		for i := range actualBytes {
			if actualBytes[i] != expectedBytes[i] {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func getMsgAndArgs(msgAndArgs ...interface{}) (string, []any) {
	var (
		msg  string
		args []any
	)

	switch len(msgAndArgs) {
	case 0:
	case 1:
		msg = msgAndArgs[0].(string)
	default:
		msg = msgAndArgs[0].(string)
		args = msgAndArgs[1:]
	}

	return msg, args
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	equal := assertEqual(expected, actual)
	if !equal {
		msg, args := getMsgAndArgs(msgAndArgs...)
		t.Errorf(msg, args...)
	}

	return equal
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	notEqual := !assertEqual(expected, actual)
	if !notEqual {
		msg, args := getMsgAndArgs(msgAndArgs...)
		t.Errorf(msg, args...)
	}

	return notEqual
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	equal := AssertEqual(t, expected, actual, msgAndArgs...)
	if !equal {
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	notEqual := AssertNotEqual(t, expected, actual, msgAndArgs...)
	if !notEqual {
		t.FailNow()
	}
}
