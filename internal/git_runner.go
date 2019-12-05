package internal

import (
	"io"
	"os/exec"
)

type GitRunner struct{}

func NewGitRunner() Runner {
	return &GitRunner{}
}

func (c *GitRunner) Run(stdout io.Writer, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = stdout
	return cmd.Run()
}
