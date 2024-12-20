package internal

import (
	"encoding/csv"
	"io"
	"os"
)

type MockRunner struct {
	echo   io.Writer
	output io.Reader
}

func NewMockRunner() Runner {
	return &MockRunner{echo: os.Stdout}
}

func (c *MockRunner) Run(sideEffects bool, stdout io.Writer, args ...string) error {
	w := csv.NewWriter(c.echo)
	w.Comma = ' '
	if err := w.Write(append([]string{"git"}, args...)); err != nil {
		return err
	}
	w.Flush()
	if stdout != nil && c.output != nil {
		_, err := io.Copy(stdout, c.output)
		return err
	}
	return nil
}

var _ Runner = (*MockRunner)(nil)
