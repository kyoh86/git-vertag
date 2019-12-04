package internal

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"os"
	"os/exec"
)

type GitCommand struct {
	push   bool
	runner Runner
}

func NewGitCommand(cwd string, dryRun, push bool) *GitCommand {
	c := &GitCommand{
		push: push,
	}

	c.runner = gitRunner(cwd)
	if dryRun {
		c.runner = dryRunner
	}

	return c
}

type Run func(*bytes.Buffer) error
type Runner func(args ...string) Run

func dryRunner(args ...string) Run {
	return func(*bytes.Buffer) error {
		w := csv.NewWriter(os.Stdout)
		w.Comma = ' '
		if err := w.Write(append([]string{"git"}, args...)); err != nil {
			return err
		}
		w.Flush()
		return nil
	}
}

func gitRunner(cwd string) Runner {
	return func(args ...string) Run {
		return func(stdout *bytes.Buffer) error {
			var cmd *exec.Cmd
			if cwd != "" {
				cmd = exec.Command("git", append([]string{"-C", cwd}, args...)...)
			} else {
				cmd = exec.Command("git", args...)
			}
			cmd.Stdout = stdout
			return cmd.Run()
		}
	}
}

func (c *GitCommand) RemoveTag(v *Semver) error {
	if err := c.runner("tag", "-d", v.String())(nil); err != nil {
		return err
	}

	if c.push {
		// UNDONE: remote name (not only origin)
		if err := c.runner("push", "origin", ":"+v.String())(nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *GitCommand) CreateTag(v *Semver, message []string, file string) error {
	args := []string{"tag"}
	for _, m := range message {
		args = append(args, "--message", m)
	}
	if file != "" {
		args = append(args, "--file", file)
	}

	if err := c.runner(append(args, v.String())...)(nil); err != nil {
		return err
	}

	if c.push {
		// UNDONE: remote name (not only origin)
		if err := c.runner("push", "origin", v.String())(nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *GitCommand) LatestVer(fetch bool) (*Semver, error) {
	if fetch {
		if err := c.runner("fetch", "--tags")(nil); err != nil {
			return nil, err
		}
	}
	var stdout bytes.Buffer
	if err := c.runner("tag", "-l")(&stdout); err != nil {
		return nil, err
	}

	latest := &Semver{}
	stream := bufio.NewScanner(&stdout)
	for stream.Scan() {
		ver, err := ParseSemver(stream.Text())
		if err != nil {
			continue
		}
		latest = GreaterSemver(latest, ver)
	}

	return latest, nil
}
