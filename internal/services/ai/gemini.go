package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Generator menghasilkan testcase dari deskripsi soal.
type Generator interface {
	GenerateTestCases(ctx context.Context, problem string, count int) ([]GeneratedTestCase, error)
}

// GeminiGenerator memanggil Gemini API (REST) dengan structured JSON output.
type GeminiGenerator struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewGeminiGenerator(apiKey, model string) *GeminiGenerator {
	return &GeminiGenerator{
		apiKey:     apiKey,
		model:      model,
		baseURL:    "https://generativelanguage.googleapis.com/v1beta",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func buildPrompt(problem string, count int) string {
	var b strings.Builder
	b.WriteString("You are generating test cases for an automated grader of C programming exercises.\n")
	b.WriteString("Problem statement:\n---\n")
	b.WriteString(problem)
	b.WriteString("\n---\n")
	fmt.Fprintf(&b, "Generate exactly %d diverse test cases (include edge cases). ", count)
	b.WriteString("Respond ONLY with a JSON array. Each element must be an object with exactly two string fields: ")
	b.WriteString(`"stdin" (the exact text fed to standard input, empty string if the program reads nothing) `)
	b.WriteString(`and "expected_output" (the exact text the correct program prints to standard output). `)
	b.WriteString("Use real newline characters inside JSON strings (JSON escape \\n), never doubled escapes. ")
	b.WriteString("Do not include explanations or markdown.")
	return b.String()
}

type generateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (g *GeminiGenerator) GenerateTestCases(ctx context.Context, problem string, count int) ([]GeneratedTestCase, error) {
	reqBody := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]string{"text": buildPrompt(problem, count)},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":      0.4,
			"responseMimeType": "application/json",
			"responseSchema": map[string]interface{}{
				"type": "ARRAY",
				"items": map[string]interface{}{
					"type": "OBJECT",
					"properties": map[string]interface{}{
						"stdin":           map[string]string{"type": "STRING"},
						"expected_output": map[string]string{"type": "STRING"},
					},
					"required": []string{"stdin", "expected_output"},
				},
			},
		},
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent", g.baseURL, g.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: gemini returned HTTP %d", ErrUnavailable, resp.StatusCode)
	}

	var out generateResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&out); err != nil {
		return nil, fmt.Errorf("%w: invalid gemini response: %v", ErrUnavailable, err)
	}
	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("%w: gemini returned no content", ErrInvalidOutput)
	}

	var text strings.Builder
	for _, part := range out.Candidates[0].Content.Parts {
		text.WriteString(part.Text)
	}
	return ParseTestCases(text.String(), count)
}
