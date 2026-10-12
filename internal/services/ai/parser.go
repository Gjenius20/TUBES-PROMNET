package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrUnavailable   = errors.New("ai service unavailable")
	ErrInvalidOutput = errors.New("ai returned invalid test cases")
)

const maxFieldBytes = 10000

// GeneratedTestCase adalah satu testcase hasil AI.
type GeneratedTestCase struct {
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

// ParseTestCases memvalidasi output AI: harus JSON array dengan schema ketat
// ([{stdin: string, expected_output: string}]) sebelum boleh disimpan ke database.
// Jika jumlah item melebihi maxCount, kelebihannya dibuang.
func ParseTestCases(raw string, maxCount int) ([]GeneratedTestCase, error) {
	cleaned := stripCodeFence(strings.TrimSpace(raw))
	if cleaned == "" {
		return nil, fmt.Errorf("%w: empty response", ErrInvalidOutput)
	}

	dec := json.NewDecoder(strings.NewReader(cleaned))
	dec.DisallowUnknownFields()

	var items []GeneratedTestCase
	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOutput, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("%w: unexpected trailing data", ErrInvalidOutput)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: no test cases returned", ErrInvalidOutput)
	}
	if maxCount > 0 && len(items) > maxCount {
		items = items[:maxCount]
	}

	for i := range items {
		items[i].Stdin = normalizeStdin(items[i].Stdin)
		items[i].ExpectedOutput = strings.TrimRight(normalizeText(items[i].ExpectedOutput), " \t\r\n")
		if items[i].ExpectedOutput == "" {
			return nil, fmt.Errorf("%w: test case %d has empty expected_output", ErrInvalidOutput, i+1)
		}
		if len(items[i].Stdin) > maxFieldBytes || len(items[i].ExpectedOutput) > maxFieldBytes {
			return nil, fmt.Errorf("%w: test case %d is too large", ErrInvalidOutput, i+1)
		}
	}
	return items, nil
}

func stripCodeFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}
	if nl := strings.Index(s, "\n"); nl >= 0 {
		s = s[nl+1:]
	}
	s = strings.TrimSpace(s)
	return strings.TrimSpace(strings.TrimSuffix(s, "```"))
}

// normalizeText menyeragamkan CRLF dan memperbaiki urutan escape ganda: bila string
// tidak memuat newline asli tetapi memuat literal "\n" (backslash + n), literal itu
// dikonversi menjadi newline sungguhan.
func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if !strings.Contains(s, "\n") && strings.Contains(s, `\n`) {
		s = strings.ReplaceAll(s, `\r\n`, "\n")
		s = strings.ReplaceAll(s, `\n`, "\n")
		s = strings.ReplaceAll(s, `\t`, "\t")
	}
	return s
}

func normalizeStdin(s string) string {
	s = normalizeText(s)
	if s != "" && !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}
