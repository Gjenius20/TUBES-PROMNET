package piston

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

var ErrUnavailable = errors.New("code execution service unavailable")

// Status ID constants (compatibility dengan Judge0).
const (
	StatusAccepted         = 3
	StatusTimeLimit        = 5
	StatusCompilationError = 6
)

// Result adalah hasil eksekusi satu program di sandbox.
type Result struct {
	Stdout            string
	Stderr            string
	CompileOutput     string
	Message           string
	StatusID          int
	StatusDescription string
}

// Runner mengabstraksi eksekusi kode sehingga mudah di-mock pada test.
type Runner interface {
	Run(ctx context.Context, sourceCode, stdin string) (*Result, error)
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type pistonFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type pistonRequest struct {
	Language           string        `json:"language"`
	Version            string        `json:"version"`
	Files              []pistonFile  `json:"files"`
	Stdin              string        `json:"stdin,omitempty"`
	CompileTimeout     int           `json:"compile_timeout"`
	RunTimeout         int           `json:"run_timeout"`
	CompileMemoryLimit int           `json:"compile_memory_limit"`
	RunMemoryLimit     int           `json:"run_memory_limit"`
}

type pistonRunOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Code   int    `json:"code"`
	Signal *string `json:"signal"`
}

type pistonCompileOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Code   int    `json:"code"`
	Signal *string `json:"signal"`
}

type pistonResponse struct {
	Language string              `json:"language"`
	Version  string              `json:"version"`
	Run      pistonRunOutput     `json:"run"`
	Compile  pistonCompileOutput `json:"compile"`
}

// Run mengirim kode ke Piston API dan mengembalikan hasilnya.
func (c *Client) Run(ctx context.Context, sourceCode, stdin string) (*Result, error) {
	payload := pistonRequest{
		Language: "c",
		Version:  "10.2.0",
		Files: []pistonFile{
			{
				Name:    "main.c",
				Content: sourceCode,
			},
		},
		Stdin:              stdin,
		CompileTimeout:     10000,
		RunTimeout:         3000,
		CompileMemoryLimit: -1,
		RunMemoryLimit:     -1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal piston request: %w", err)
	}

	url := c.baseURL + "/execute"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build piston request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: unexpected HTTP status %d", ErrUnavailable, resp.StatusCode)
	}

	var out pistonResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&out); err != nil {
		return nil, fmt.Errorf("%w: invalid response: %v", ErrUnavailable, err)
	}

	// Map Piston response to Result.
	// If compilation failed, statusID = 6 (CompilationError).
	statusID := 3 // Default: Accepted
	statusDesc := "Accepted"
	compileOutput := out.Compile.Stderr
	stdout := out.Run.Stdout
	stderr := out.Run.Stderr

	if out.Compile.Code != 0 {
		statusID = 6 // CompilationError
		statusDesc = "Compilation Error"
		stdout = ""
		stderr = out.Compile.Stderr
	} else if out.Run.Code != 0 {
		statusID = 5 // TimeLimit (generic non-zero exit)
		statusDesc = "Runtime Error"
	}

	return &Result{
		Stdout:            stdout,
		Stderr:            stderr,
		CompileOutput:     compileOutput,
		Message:           "",
		StatusID:          statusID,
		StatusDescription: statusDesc,
	}, nil
}
