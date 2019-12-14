package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* SPEC: https://semver.org/ */

func TestSemver(t *testing.T) {

	t.Run("parse and stringify", func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			input := "v1.2.3-alpha.4+build.5"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("build", func(t *testing.T) {
			input := "v1.2.3+build.5"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("pre-release", func(t *testing.T) {
			input := "v1.2.3-alpha.4"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("patch", func(t *testing.T) {
			input := "v1.2.3"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("minor", func(t *testing.T) {
			input := "v1.2-alpha.4"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, "v1.2.0-alpha.4", semver.String())
		})

		t.Run("minor", func(t *testing.T) {
			input := "v1.2"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, "v1.2.0", semver.String())
		})

		t.Run("major with note", func(t *testing.T) {
			input := "v1-alpha.4"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, "v1.0.0-alpha.4", semver.String())
		})

		t.Run("major", func(t *testing.T) {
			input := "v1"
			semver, err := ParseSemver(input)
			require.NoError(t, err)
			assert.Equal(t, "v1.0.0", semver.String())
		})

	})

	t.Run("manipulate", func(t *testing.T) {
		t.Run("patch", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseSemver("v1.2.3-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v1.2.4", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseSemver("v1.2.3-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v1.3.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseSemver("v1.2.3-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v2.0.0", semver.Update().Major().Apply().String())
			})
		})

		t.Run("minor", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseSemver("v1.2-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v1.2.1", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseSemver("v1.2-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v1.3.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseSemver("v1.2-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v2.0.0", semver.Update().Major().Apply().String())
			})
		})

		t.Run("major", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseSemver("v1-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v1.0.1", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseSemver("v1-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v1.1.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseSemver("v1-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "v2.0.0", semver.Update().Major().Apply().String())
			})
		})

		t.Run("comparison", func(t *testing.T) {
			t.Run("patch > patch", func(t *testing.T) {})
			t.Run("patch = patch", func(t *testing.T) {})
			t.Run("patch < patch", func(t *testing.T) {})

			t.Run("patch > minor", func(t *testing.T) {})
			t.Run("patch < minor", func(t *testing.T) {})

			t.Run("patch > major", func(t *testing.T) {})
			t.Run("patch < major", func(t *testing.T) {})
		})
	})
}
