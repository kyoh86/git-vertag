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
			input := "1.2.3-alpha.4+build.5"
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("build", func(t *testing.T) {
			input := "1.2.3+build.5"
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("pre-release", func(t *testing.T) {
			input := "1.2.3-alpha.4"
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("patch", func(t *testing.T) {
			input := "1.2.3"
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, input, semver.String())
		})

		t.Run("minor-prerelease+build", func(t *testing.T) {
			input := "v1.2-alpha.4+build.5"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.2.0-alpha.4+build.5", semver.String())
		})

		t.Run("minor-prerelease", func(t *testing.T) {
			input := "v1.2-alpha.4"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.2.0-alpha.4", semver.String())
		})

		t.Run("minor+build", func(t *testing.T) {
			input := "v1.2+build.5"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.2.0+build.5", semver.String())
		})

		t.Run("minor", func(t *testing.T) {
			input := "v1.2"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.2.0", semver.String())
		})

		t.Run("major-prerelease+build", func(t *testing.T) {
			input := "v1-alpha.4+build.5"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.0.0-alpha.4+build.5", semver.String())
		})

		t.Run("major-prerelease", func(t *testing.T) {
			input := "v1-alpha.4"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.0.0-alpha.4", semver.String())
		})

		t.Run("major+build", func(t *testing.T) {
			input := "v1+build.5"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.0.0+build.5", semver.String())
		})

		t.Run("major", func(t *testing.T) {
			input := "v1"
			_, err := Parse(input)
			assert.Error(t, err)
			semver, err := ParseTolerant(input)
			require.NoError(t, err)
			assert.Equal(t, "1.0.0", semver.String())
		})

	})

	t.Run("manipulate", func(t *testing.T) {
		t.Run("source: pre-release+build", func(t *testing.T) {
			t.Run("set build", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4+build.5")
				require.NoError(t, err)
				v := semver.Update().
					Build(MustParseBuildID("build"), MustParseBuildID("6")).Apply()
				assert.Equal(t, "1.2.3-alpha.4+build.6", v.String())
			})
			t.Run("set pre-release and build", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4+build.5")
				require.NoError(t, err)
				v := semver.Update().
					PreRelease(MustParsePreReleaseID("beta"), MustParsePreReleaseID("6")).
					Build(MustParseBuildID("build"), MustParseBuildID("7")).Apply()
				assert.Equal(t, "1.2.3-beta.6+build.7", v.String())
			})
			t.Run("set pre-release", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4+build.5")
				require.NoError(t, err)
				v := semver.Update().
					PreRelease(MustParsePreReleaseID("beta"), MustParsePreReleaseID("6")).Apply()
				assert.Equal(t, "1.2.3-beta.6", v.String())
			})
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4+build.5")
				require.NoError(t, err)
				assert.Equal(t, "1.2.4", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4+build.5")
				require.NoError(t, err)
				assert.Equal(t, "1.3.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4+build.5")
				require.NoError(t, err)
				assert.Equal(t, "2.0.0", semver.Update().Major().Apply().String())
			})
		})

		t.Run("source: pre-release", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "1.2.4", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "1.3.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2.3-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "2.0.0", semver.Update().Major().Apply().String())
			})
		})

		t.Run("minor", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "1.2.1", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "1.3.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseTolerant("v1.2-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "2.0.0", semver.Update().Major().Apply().String())
			})
		})

		t.Run("major", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := ParseTolerant("v1-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "1.0.1", semver.Update().Patch().Apply().String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := ParseTolerant("v1-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "1.1.0", semver.Update().Minor().Apply().String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := ParseTolerant("v1-alpha.4")
				require.NoError(t, err)
				assert.Equal(t, "2.0.0", semver.Update().Major().Apply().String())
			})
		})
	})
}
