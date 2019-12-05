package internal

import (
	"io"
)

type TagCommand interface {
	Run(stdout io.Writer, args ...string) error
}
