package autograder

import "testing"

func TestMatch(t *testing.T) {
	cases := []struct {
		a, e string
		want bool
	}{
		{"Hello World\n", "hello   world", true},
		{"3.1415926", "3.14159265", true},
		{"3.14", "3.15", false},
		{"1 2 3", "1 2", false},
		{"", "", true},
	}
	for _, c := range cases {
		if got := Match(c.a, c.e); got != c.want {
			t.Errorf("Match(%q,%q)=%v want %v", c.a, c.e, got, c.want)
		}
	}
}
