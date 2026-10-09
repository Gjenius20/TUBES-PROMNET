package ai

import (
	"errors"
	"testing"
)

func TestParseTestCasesValid(t *testing.T) {
	raw := `[{"stdin":"1 2","expected_output":"3\n"},{"stdin":"","expected_output":"hello"}]`
	got, err := ParseTestCases(raw, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].Stdin != "1 2\n" || got[0].ExpectedOutput != "3" {
		t.Fatalf("unexpected normalization: %+v", got[0])
	}
	if got[1].Stdin != "" {
		t.Fatalf("empty stdin must stay empty, got %q", got[1].Stdin)
	}
}

func TestParseTestCasesDoubleEscaped(t *testing.T) {
	raw := "[{\"stdin\":\"1\\\\n2\",\"expected_output\":\"3\"}]"
	got, err := ParseTestCases(raw, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0].Stdin != "1\n2\n" {
		t.Fatalf("double escaped newline not fixed: %q", got[0].Stdin)
	}
}

func TestParseTestCasesCodeFenceAndTruncate(t *testing.T) {
	raw := "```json\n[{\"stdin\":\"a\",\"expected_output\":\"b\"},{\"stdin\":\"c\",\"expected_output\":\"d\"}]\n```"
	got, err := ParseTestCases(raw, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected truncation to 1 item, got %d", len(got))
	}
}

func TestParseTestCasesInvalid(t *testing.T) {
	cases := map[string]string{
		"not json":      "hello",
		"object":        `{"stdin":"a","expected_output":"b"}`,
		"empty array":   `[]`,
		"empty output":  `[{"stdin":"a","expected_output":"  "}]`,
		"unknown field": `[{"stdin":"a","expected_output":"b","x":1}]`,
		"wrong type":    `[{"stdin":1,"expected_output":"b"}]`,
		"trailing":      `[{"stdin":"a","expected_output":"b"}] extra`,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseTestCases(raw, 5); !errors.Is(err, ErrInvalidOutput) {
				t.Fatalf("expected ErrInvalidOutput, got %v", err)
			}
		})
	}
}
