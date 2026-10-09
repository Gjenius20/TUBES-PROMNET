package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	key, model string
	http       *http.Client
}

func New(key, model string) *Client {
	return &Client{key: key, model: model, http: &http.Client{Timeout: 60 * time.Second}}
}

type TestCase struct {
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

const promptTpl = `You write test cases for a C programming exercise.
Problem statement:
%s

Return ONLY a JSON array of exactly %d objects, each with the string fields "stdin" (program input) and "expected_output" (exact expected stdout). No prose, no markdown.`

func (c *Client) GenerateTestCases(ctx context.Context, statement string, n int) ([]TestCase, error) {
	payload, _ := json.Marshal(map[string]any{
		"contents":         []any{map[string]any{"parts": []any{map[string]string{"text": fmt.Sprintf(promptTpl, statement, n)}}}},
		"generationConfig": map[string]any{"responseMimeType": "application/json"},
	})

	modelName := strings.TrimPrefix(c.model, "models/")
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}
	url := "https://generativelanguage.googleapis.com/v1beta/models/" + modelName + ":generateContent"

	var resp *http.Response
	var err error
	maxRetries := 4
	backoff := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if reqErr != nil {
			return nil, reqErr
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", c.key)

		resp, err = c.http.Do(req)
		if err != nil {
			// Network-level error, wait and retry
			if attempt == maxRetries {
				return nil, fmt.Errorf("network error after %d retries: %w", maxRetries, err)
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		// Retry on 503 (Service Unavailable) or 429 (Rate Limit / High Demand)
		if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			if attempt == maxRetries {
				break
			}
			time.Sleep(backoff)
			backoff *= 2 // Wait 2s, 4s, 8s...
			continue
		}

		// Stop retrying on any other status code (200, 400, 401, 404, etc.)
		break
	}

	if resp == nil {
		return nil, errors.New("failed to get a response from gemini after retries")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini returned HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var out struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("gemini returned no content")
	}
	return ParseTestCases(out.Candidates[0].Content.Parts[0].Text)
}

// ParseTestCases validates the model output against the strict schema.
func ParseTestCases(raw string) ([]TestCase, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var tcs []TestCase
	if err := json.Unmarshal([]byte(raw), &tcs); err != nil {
		return nil, fmt.Errorf("ai output is not a valid testcase array: %w", err)
	}
	if len(tcs) == 0 {
		return nil, errors.New("ai returned an empty testcase array")
	}
	for i := range tcs {
		tcs[i].Stdin = unescape(tcs[i].Stdin)
		tcs[i].ExpectedOutput = unescape(tcs[i].ExpectedOutput)
		if strings.TrimSpace(tcs[i].ExpectedOutput) == "" {
			return nil, fmt.Errorf("testcase %d has empty expected_output", i+1)
		}
	}
	return tcs, nil
}

func unescape(s string) string {
	if !strings.Contains(s, "\n") {
		s = strings.ReplaceAll(s, `\n`, "\n")
	}
	return s
}