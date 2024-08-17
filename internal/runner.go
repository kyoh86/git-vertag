package internal

import (
	"io"
)

type Runner interface {
	Run(sideEffects bool, stdout io.Writer, args ...string) error
}
