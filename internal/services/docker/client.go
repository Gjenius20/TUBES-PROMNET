package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var ErrUnavailable = errors.New("docker execution service unavailable")

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

// Config mengatur batasan sandbox. Nilai nol/string kosong akan diisi default.
type Config struct {
	ImageName string
	Timeout   time.Duration
	MemoryMB  int
	CPULimit  string
}

// Client menjalankan kode C di container Docker (gcc) yang terisolasi.
type Client struct {
	cfg Config
}

func NewClient(cfg Config) *Client {
	if cfg.ImageName == "" {
		cfg.ImageName = "gcc:12"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.MemoryMB <= 0 {
		cfg.MemoryMB = 128
	}
	if cfg.CPULimit == "" {
		cfg.CPULimit = "0.5"
	}
	return &Client{cfg: cfg}
}

func (c *Client) Run(ctx context.Context, sourceCode, stdin string) (*Result, error) {
	runCtx := ctx
	cancel := func() {}
	if c.cfg.Timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, c.cfg.Timeout)
	}
	defer cancel()

	// Source di-encode base64 agar aman dilewatkan ke shell.
	encodedSource := base64.StdEncoding.EncodeToString([]byte(sourceCode))
	// Gunakan printf, bukan echo, agar escape sequence tidak ditafsirkan shell.
	script := fmt.Sprintf(`printf '%s' | base64 -d > main.c && gcc main.c -o main && ./main`, encodedSource)

	// Sandbox: tanpa network, dengan batas memori/CPU/pid untuk mencegah DoS.
	args := []string{"run", "--rm", "-i", "--network", "none"}
	if c.cfg.MemoryMB > 0 {
		args = append(args, "--memory", fmt.Sprintf("%dm", c.cfg.MemoryMB))
	}
	if c.cfg.CPULimit != "" {
		args = append(args, "--cpus", c.cfg.CPULimit)
	}
	args = append(args, "--pids-limit", "64", c.cfg.ImageName, "sh", "-c", script)

	cmd := exec.CommandContext(runCtx, "docker", args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	res := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	stderrStr := stderr.String()

	switch {
	case isCompilationError(stderrStr):
		res.StatusID = StatusCompilationError
		res.StatusDescription = "Compilation Error"
		res.CompileOutput = stderrStr
	case runCtx.Err() == context.DeadlineExceeded:
		res.StatusID = StatusTimeLimit
		res.StatusDescription = "Time Limit Exceeded"
	case err != nil:
		res.StatusID = StatusTimeLimit
		res.StatusDescription = "Runtime Error"
	default:
		res.StatusID = StatusAccepted
		res.StatusDescription = "Accepted"
	}

	return res, nil
}

// isCompilationError mendeteksi error kompilasi gcc dari stderr.
func isCompilationError(stderr string) bool {
	return strings.Contains(stderr, "error:") ||
		strings.Contains(stderr, "Error:") ||
		strings.Contains(stderr, "fatal error:") ||
		strings.Contains(stderr, "undefined reference") ||
		strings.Contains(stderr, "ld returned")
}
