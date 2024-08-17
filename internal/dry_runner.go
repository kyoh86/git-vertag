package internal

import (
	"encoding/csv"
	"io"
	"os"
)

type DryRunner struct {
	Runner
	echo *csv.Writer
}

func NewDryRunner() Runner {
	w := csv.NewWriter(os.Stdout)
	w.Comma = ' '
	return &DryRunner{
		Runner: NewGitRunner(),
		echo:   w,
	}
}

func (c *DryRunner) Run(sideEffects bool, stdout io.Writer, args ...string) error {
	if !sideEffects {
		return c.Runner.Run(sideEffects, stdout, args...)
	}
	if err := c.echo.Write(append([]string{"may", "calls:", "git"}, args...)); err != nil {
		return err
	}
	c.echo.Flush()
	return nil
}

var _ Runner = (*DryRunner)(nil)
