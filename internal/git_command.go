package internal

import (
	"io"
	"os/exec"
)

type GitCommand struct{}

func NewGitCommand() TagCommand {
	return &GitCommand{}
}

func (c *GitCommand) Run(stdout io.Writer, args ...string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("git", args...)
	cmd.Stdout = stdout
	return cmd.Run()
}
