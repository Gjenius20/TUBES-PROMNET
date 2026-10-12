package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrUnavailable = errors.New("docker execution service unavailable")

const (
	StatusAccepted         = 3
	StatusTimeLimit        = 5
	StatusCompilationError = 6
)

type Result struct {
	Stdout            string
	Stderr            string
	CompileOutput     string
	Message           string
	StatusID          int
	StatusDescription string
}

type Runner interface {
	Run(ctx context.Context, sourceCode, stdin string) (*Result, error)
}

type Client struct {
	imageName string
}

func NewClient(imageName string) *Client {
	if imageName == "" {
		imageName = "gcc:12"
	}
	return &Client{imageName: imageName}
}

func (c *Client) Run(ctx context.Context, sourceCode, stdin string) (*Result, error) {
	// Script: create file, compile, run
	script := fmt.Sprintf(`cat > main.c << 'EOF'
%s
EOF
gcc main.c -o main && ./main`, sourceCode)

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i",
		c.imageName,
		"sh", "-c", script,
	)

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

	if err != nil {
		// Identify compilation error
		if strings.Contains(stderr.String(), "error:") {
			res.StatusID = StatusCompilationError
			res.StatusDescription = "Compilation Error"
			res.CompileOutput = stderr.String()
		} else {
			res.StatusID = StatusTimeLimit
			res.StatusDescription = "Runtime Error"
		}
	} else {
		res.StatusID = StatusAccepted
		res.StatusDescription = "Accepted"
	}

	return res, nil
}
