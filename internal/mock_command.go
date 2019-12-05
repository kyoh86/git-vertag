package internal

import (
	"encoding/csv"
	"io"
	"os"
)

type MockCommand struct {
	echo   io.Writer
	output io.Reader
}

func NewMockCommand() TagCommand {
	return &MockCommand{echo: os.Stdout}
}

func (c *MockCommand) Run(stdout io.Writer, args ...string) error {
	w := csv.NewWriter(c.echo)
	w.Comma = ' '
	if err := w.Write(append([]string{"git"}, args...)); err != nil {
		return err
	}
	w.Flush()
	if stdout != nil && c.output != nil {
		io.Copy(stdout, c.output)
	}
	return nil
}
