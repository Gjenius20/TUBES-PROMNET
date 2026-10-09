package autograder

import (
	"math"
	"strconv"
	"strings"
)

const floatTolerance = 1e-6

// Match compares program output with the expected output token by token.
// It ignores whitespace layout and letter case, and accepts floats within a small tolerance.
func Match(actual, expected string) bool {
	a, e := strings.Fields(actual), strings.Fields(expected)
	if len(a) != len(e) {
		return false
	}
	for i := range a {
		if strings.EqualFold(a[i], e[i]) {
			continue
		}
		fa, errA := strconv.ParseFloat(a[i], 64)
		fe, errE := strconv.ParseFloat(e[i], 64)
		if errA != nil || errE != nil || math.Abs(fa-fe) > floatTolerance {
			return false
		}
	}
	return true
}
