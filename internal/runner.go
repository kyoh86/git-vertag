package internal

import (
	"io"
)

type Runner interface {
	Run(stdout io.Writer, args ...string) error
}
