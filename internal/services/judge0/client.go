package judge0

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// languageIDC adalah ID bahasa "C (GCC 9.2.0)" pada Judge0 CE 1.13.x.
const languageIDC = 50

// Status ID Judge0 yang relevan.
const (
	StatusAccepted         = 3
	StatusTimeLimit        = 5
	StatusCompilationError = 6
)

var ErrUnavailable = errors.New("judge0 execution service unavailable")

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

type submissionRequest struct {
	LanguageID    int     `json:"language_id"`
	SourceCode    string  `json:"source_code"`
	Stdin         string  `json:"stdin,omitempty"`
	CPUTimeLimit  float64 `json:"cpu_time_limit"`
	WallTimeLimit float64 `json:"wall_time_limit"`
}

type submissionResponse struct {
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

// Run mengirim kode ke Judge0 (sinkron, wait=true) dan mengembalikan hasilnya.
// Semua kegagalan jaringan / status HTTP tak terduga dibungkus ErrUnavailable.
func (c *Client) Run(ctx context.Context, sourceCode, stdin string) (*Result, error) {
	payload := submissionRequest{
		LanguageID:    languageIDC,
		SourceCode:    base64.StdEncoding.EncodeToString([]byte(sourceCode)),
		CPUTimeLimit:  2,
		WallTimeLimit: 10,
	}
	if stdin != "" {
		payload.Stdin = base64.StdEncoding.EncodeToString([]byte(stdin))
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal judge0 request: %w", err)
	}

	url := c.baseURL + "/submissions?base64_encoded=true&wait=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build judge0 request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: unexpected HTTP status %d", ErrUnavailable, resp.StatusCode)
	}

	var out submissionResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&out); err != nil {
		return nil, fmt.Errorf("%w: invalid response: %v", ErrUnavailable, err)
	}

	return &Result{
		Stdout:            decode(out.Stdout),
		Stderr:            decode(out.Stderr),
		CompileOutput:     decode(out.CompileOutput),
		Message:           decode(out.Message),
		StatusID:          out.Status.ID,
		StatusDescription: out.Status.Description,
	}, nil
}

// decode membaca field base64 dari Judge0 (decoder Go mengabaikan newline).
func decode(p *string) string {
	if p == nil {
		return ""
	}
	b, err := base64.StdEncoding.DecodeString(*p)
	if err != nil {
		return *p
	}
	return string(b)
}
