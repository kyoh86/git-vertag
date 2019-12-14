package semver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	// ORDERS from SPEC.
	// see: https://semver.org/
	orders := [][]string{
		{"1.0.0", "2.0.0", "2.1.0", "2.1.1"},
		{"1.0.0-alpha", "1.0.0"},
		{"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta", "1.0.0-beta", "1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0"},
	}

	for _, order := range orders {
		for i1, item1 := range order {
			for i2, item2 := range order {
				t.Run(fmt.Sprintf("compare %s and %s", item1, item2), func(t *testing.T) {
					v1, err := Parse(item1)
					require.NoError(t, err)
					v2, err := Parse(item2)
					require.NoError(t, err)

					expect := 0
					switch {
					case i1 < i2:
						expect = -1
					case i1 > i2:
						expect = 1
					}

					assert.Equal(t, expect, Compare(v1, v2))
				})
			}
		}
	}
}
