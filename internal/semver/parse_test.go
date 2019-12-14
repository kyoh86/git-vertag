package semver

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		source   string
		expected Semver
	}{
		{"1.0.0", Semver{Major: 1, Minor: 0, Patch: 0}},
		{"1.10.0-alpha", Semver{Major: 1, Minor: 10, Patch: 0, PreRelease: PreRelease{{str: "alpha"}}}},
		{"1.0.0-alpha.1", Semver{Major: 1, Minor: 0, Patch: 0, PreRelease: PreRelease{{str: "alpha"}, {str: "1", num: 1, isNum: true}}}},
		{"1.0.0-alpha.beta", Semver{Major: 1, Minor: 0, Patch: 0, PreRelease: PreRelease{{str: "alpha"}, {str: "beta"}}}},
		{"1.0.0-beta.11", Semver{Major: 1, Minor: 0, Patch: 0, PreRelease: PreRelease{{str: "beta"}, {str: "11", num: 11, isNum: true}}}},
		{"1.10.0+alpha", Semver{Major: 1, Minor: 10, Patch: 0, Build: Build{"alpha"}}},
		{"1.0.0+alpha.1", Semver{Major: 1, Minor: 0, Patch: 0, Build: Build{"alpha", "1"}}},
		{"1.0.0+alpha.beta", Semver{Major: 1, Minor: 0, Patch: 0, Build: Build{"alpha", "beta"}}},
		{"1.0.0+beta.11", Semver{Major: 1, Minor: 0, Patch: 0, Build: Build{"beta", "11"}}},
		{"1.10.0-alpha+beta", Semver{Major: 1, Minor: 10, Patch: 0, PreRelease: PreRelease{{str: "alpha"}}, Build: Build{"beta"}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			actual, err := Parse(tt.source)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
