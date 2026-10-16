package docker

import "testing"

func TestIsCompilationError(t *testing.T) {
	cases := []struct {
		name   string
		stderr string
		want   bool
	}{
		{"syntax error", "main.c:3:5: error: expected ';' before '}'", true},
		{"missing header", "main.c:1:10: fatal error: stdio.h: No such file or directory", true},
		{"undefined reference", "/usr/bin/ld: /tmp/ccX.o: undefined reference to `foo'", true},
		{"linker failed", "collect2: error: ld returned 1 exit status", true},
		{"empty stderr", "", false},
		{"segfault is runtime", "Segmentation fault (core dumped)", false},
		{"normal program output", "Hello, World!", false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCompilationError(tt.stderr); got != tt.want {
				t.Fatalf("isCompilationError(%q) = %v, want %v", tt.stderr, got, tt.want)
			}
		})
	}
}
