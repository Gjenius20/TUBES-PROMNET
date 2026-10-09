package autograder

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
		want     bool
	}{
		{"exact", "42\n", "42\n", true},
		{"whitespace differences", "1   2\r\n3 \n\n", "1 2 3", true},
		{"different order", "2 1", "1 2", false},
		{"different length", "1 2 3", "1 2", false},
		{"float within tolerance", "3.1415927", "3.14159265", true},
		{"float outside tolerance", "3.15", "3.14", false},
		{"keyword mismatch", "Yes", "yes", false},
		{"keyword match", "Result: 10", "Result: 10", true},
		{"empty both", "", "\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.actual, tt.expected); got != tt.want {
				t.Fatalf("Match(%q, %q) = %v, want %v", tt.actual, tt.expected, got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	got := Normalize("a  \r\nb\t\n\n\n")
	if got != "a\nb" {
		t.Fatalf("Normalize returned %q", got)
	}
}
