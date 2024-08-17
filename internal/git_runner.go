package internal

import (
	"io"
	"os/exec"
)

type GitRunner struct{}

func NewGitRunner() Runner {
	return &GitRunner{}
}

func (c *GitRunner) Run(sideEffects bool, stdout io.Writer, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = stdout
	return cmd.Run()
}

var _ Runner = (*GitRunner)(nil)
