package judge0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const languageC = 50 // C (GCC 9.2.0)

const StatusCompileError = 6

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 30 * time.Second}}
}

type Result struct {
	Stdout, Stderr, CompileOutput string
	StatusID                      int
	Description                   string
}

func (c *Client) Run(ctx context.Context, source, stdin string) (*Result, error) {
	body, _ := json.Marshal(map[string]any{"source_code": source, "language_id": languageC, "stdin": stdin})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/submissions?base64_encoded=false&wait=true", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("judge0 unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("judge0 returned HTTP %d", resp.StatusCode)
	}
	var out struct {
		Stdout        *string `json:"stdout"`
		Stderr        *string `json:"stderr"`
		CompileOutput *string `json:"compile_output"`
		Status        struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	val := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	return &Result{val(out.Stdout), val(out.Stderr), val(out.CompileOutput), out.Status.ID, out.Status.Description}, nil
}
