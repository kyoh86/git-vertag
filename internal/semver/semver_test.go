package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/* SPEC: https://semver.org/ */

func TestSemver(t *testing.T) {
	t.Run("parse strict and stringify", func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			input := "1.2.3-pre-release.4+build-ver.5"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("big number", func(t *testing.T) {
			input := "123.456.789-pre-release.234+build-ver.567"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("build", func(t *testing.T) {
			input := "1.2.3+build-ver.5"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("pre-release", func(t *testing.T) {
			input := "1.2.3-pre-release.4"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("patch", func(t *testing.T) {
			input := "1.2.3"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})
	})

	t.Run("parse strict error", func(t *testing.T) {
		t.Run("zero-started major", func(t *testing.T) {
			input := "01.2.3-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("zero-started minor", func(t *testing.T) {
			input := "1.02.3-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("zero-started pre-release", func(t *testing.T) {
			input := "1.2.3-pre-release.04+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("zero-started pre-release in tail", func(t *testing.T) {
			input := "1.2.3-pre-release.04"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		// NOTE: build identifiers can be zero-started

		t.Run("empty major", func(t *testing.T) {
			input := ".2.3-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty minor", func(t *testing.T) {
			input := "1..3-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty patch", func(t *testing.T) {
			input := "1.2.-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty first pre-release identifier", func(t *testing.T) {
			input := "1.2.3-.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty last pre-release identifier", func(t *testing.T) {
			input := "1.2.3-pre-release.+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty pre-release", func(t *testing.T) {
			input := "1.2.3-+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty first build identifier", func(t *testing.T) {
			input := "1.2.3-pre-release.4+.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty last build identifier", func(t *testing.T) {
			input := "1.2.3-pre-release.4+build."
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("empty build", func(t *testing.T) {
			input := "1.2.3-pre-release.4+"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("alphabetical major", func(t *testing.T) {
			input := "a.2.3-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("alphabetical minor", func(t *testing.T) {
			input := "1.a.3-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("alphabetical patch", func(t *testing.T) {
			input := "1.2.a-pre-release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("invalid char in pre-release", func(t *testing.T) {
			input := "1.2.3-pre*release.4+build-ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

		t.Run("invalid char in build", func(t *testing.T) {
			input := "1.2.3-pre-release.4+build*ver.5"
			_, err := Parse(input)
			assert.Error(t, err)
		})

	})

	t.Run("manipulate", func(t *testing.T) {
		t.Run("source: pre-release+build-ver", func(t *testing.T) {
			t.Run("set build", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4+build-ver.5")
				require.NoError(t, err)
				v, err := semver.Update().
					Build(MustParseBuildID("build-ver"), MustParseBuildID("6")).Apply()
				require.NoError(t, err)
				assert.Equal(t, "1.2.3-pre-release.4+build-ver.6", v.String())
			})
			t.Run("set pre-release and build", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4+build-ver.5")
				require.NoError(t, err)
				v, err := semver.Update().
					PreRelease(MustParsePreReleaseID("beta"), MustParsePreReleaseID("6")).
					Build(MustParseBuildID("build-ver"), MustParseBuildID("7")).Apply()
				require.NoError(t, err)
				assert.Equal(t, "1.2.3-beta.6+build-ver.7", v.String())
			})
			t.Run("set pre-release", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4+build-ver.5")
				require.NoError(t, err)
				v, err := semver.Update().
					PreRelease(MustParsePreReleaseID("beta"), MustParsePreReleaseID("6")).Apply()
				require.NoError(t, err)
				assert.Equal(t, "1.2.3-beta.6", v.String())
			})
			t.Run("increment patch", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4+build-ver.5")
				require.NoError(t, err)
				v, err := semver.Update().Patch().Apply()
				require.NoError(t, err)
				assert.Equal(t, "1.2.4", v.String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4+build-ver.5")
				require.NoError(t, err)
				v, err := semver.Update().Minor().Apply()
				assert.Equal(t, "1.3.0", v.String())
				require.NoError(t, err)
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4+build-ver.5")
				require.NoError(t, err)
				v, err := semver.Update().Major().Apply()
				require.NoError(t, err)
				assert.Equal(t, "2.0.0", v.String())
			})
		})

		t.Run("source: pre-release", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4")
				require.NoError(t, err)
				v, err := semver.Update().Patch().Apply()
				assert.Equal(t, "1.2.4", v.String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4")
				require.NoError(t, err)
				v, err := semver.Update().Minor().Apply()
				require.NoError(t, err)
				assert.Equal(t, "1.3.0", v.String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := Parse("1.2.3-pre-release.4")
				require.NoError(t, err)
				v, err := semver.Update().Major().Apply()
				require.NoError(t, err)
				assert.Equal(t, "2.0.0", v.String())
			})
		})

	})
}
