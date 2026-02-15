package tabletest

import (
	"testing"
	"time"
)

var parseTableTests = map[string]struct {
	input    string
	expected time.Duration
	hasError bool
}{
	"empty string":                  {"", 0, true},
	"zero":                          {"0", 0, false},
	"invalid duration":              {"invalid", 0, true},
	"overflow":                      {"1000000000000000000000", 0, true},
	"overflow2 ":                    {"9223372036854775808", 0, true},
	"fractional overflow":           {".1111111111111111111111", 0, true},
	"fractional overflow 2":         {".9223372036854775808", 0, true},
	"invalid duration fraction":     {".", 0, true},
	"missing unit":                  {"42", 0, true},
	"unknown unit":                  {"42.parrots", 0, true},
	"overflow with unit":            {"9223372036854775807us", 0, true},
	"overflow fraction unit":        {"9223372036854775807.9223372036854775807h", 0, true},
	"fractional overflow after add": {"2562047.9h", 0, true},
	"success negative duration":     {"-4.2h", -4*3600*1000000000 - 12*60*1000000000, false},
}

func TestParseDuration(t *testing.T) {
	for name, testcase := range parseTableTests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, err := ParseDuration(testcase.input)
			if testcase.hasError && err == nil {
				t.Errorf("ParseDuration didn't return an error")
			}
			if res != testcase.expected {
				t.Errorf("Got: %v\nWant: %v\n", res, testcase.expected)
			}
		})
	}
}
