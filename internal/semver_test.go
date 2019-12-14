package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemver(t *testing.T) {

	t.Run("parse and stringify", func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			input := "v1.2.3-notes"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, LevelPatch, semver.Level())
			assert.Equal(t, input, semver.String())
			assert.Equal(t, "v1.2.3", semver.PatchString())
			assert.Equal(t, "v1.2", semver.MinorString())
			assert.Equal(t, "v1", semver.MajorString())
		})

		t.Run("full", func(t *testing.T) {
			input := "v1.2.3"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, LevelPatch, semver.Level())
			assert.Equal(t, input, semver.String())
			assert.Equal(t, "v1.2.3", semver.PatchString())
			assert.Equal(t, "v1.2", semver.MinorString())
			assert.Equal(t, "v1", semver.MajorString())
		})

		t.Run("minor with note", func(t *testing.T) {
			input := "v1.2-notes"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, LevelMinor, semver.Level())
			assert.Equal(t, input, semver.String())
			assert.Equal(t, "v1.2.0", semver.PatchString())
			assert.Equal(t, "v1.2", semver.MinorString())
			assert.Equal(t, "v1", semver.MajorString())
		})

		t.Run("minor", func(t *testing.T) {
			input := "v1.2"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, LevelMinor, semver.Level())
			assert.Equal(t, input, semver.String())
			assert.Equal(t, "v1.2.0", semver.PatchString())
			assert.Equal(t, "v1.2", semver.MinorString())
			assert.Equal(t, "v1", semver.MajorString())
		})

		t.Run("major with note", func(t *testing.T) {
			input := "v1-notes"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, LevelMajor, semver.Level())
			assert.Equal(t, input, semver.String())
			assert.Equal(t, "v1.0.0", semver.PatchString())
			assert.Equal(t, "v1.0", semver.MinorString())
			assert.Equal(t, "v1", semver.MajorString())
		})

		t.Run("major", func(t *testing.T) {
			input := "v1"
			semver, err := Parse(input)
			require.NoError(t, err)
			assert.Equal(t, LevelMajor, semver.Level())
			assert.Equal(t, input, semver.String())
			assert.Equal(t, "v1.0.0", semver.PatchString())
			assert.Equal(t, "v1.0", semver.MinorString())
			assert.Equal(t, "v1", semver.MajorString())
		})

	})

	t.Run("manipulate", func(t *testing.T) {
		t.Run("patch", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := Parse("v1.2.3-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.2.4", semver.Increment(LevelPatch).String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := Parse("v1.2.3-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.3.0", semver.Increment(LevelMinor).String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := Parse("v1.2.3-notes")
				require.NoError(t, err)
				assert.Equal(t, "v2.0.0", semver.Increment(LevelMajor).String())
			})
			t.Run("decrement patch", func(t *testing.T) {
				semver, err := Parse("v1.2.3-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.2.2", semver.Decrement(LevelPatch).String())
			})
			t.Run("decrement minor", func(t *testing.T) {
				semver, err := Parse("v1.2.3-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.1.0", semver.Decrement(LevelMinor).String())
			})
			t.Run("decrement major", func(t *testing.T) {
				semver, err := Parse("v1.2.3-notes")
				require.NoError(t, err)
				assert.Equal(t, "v0.0.0", semver.Decrement(LevelMajor).String())
			})
		})

		t.Run("minor", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := Parse("v1.2-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.2.1", semver.Increment(LevelPatch).String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := Parse("v1.2-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.3", semver.Increment(LevelMinor).String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := Parse("v1.2-notes")
				require.NoError(t, err)
				assert.Equal(t, "v2.0", semver.Increment(LevelMajor).String())
			})
			t.Run("decrement patch", func(t *testing.T) {
				semver, err := Parse("v1.2-notes")
				require.NoError(t, err)
				assert.PanicsWithValue(t, "undefined level", func() {
					semver.Decrement(LevelPatch)
				})
			})
			t.Run("decrement minor", func(t *testing.T) {
				semver, err := Parse("v1.2-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.1", semver.Decrement(LevelMinor).String())
			})
			t.Run("decrement major", func(t *testing.T) {
				semver, err := Parse("v1.2-notes")
				require.NoError(t, err)
				assert.Equal(t, "v0.0", semver.Decrement(LevelMajor).String())
			})
		})

		t.Run("major", func(t *testing.T) {
			t.Run("increment patch", func(t *testing.T) {
				semver, err := Parse("v1-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.0.1", semver.Increment(LevelPatch).String())
			})
			t.Run("increment minor", func(t *testing.T) {
				semver, err := Parse("v1-notes")
				require.NoError(t, err)
				assert.Equal(t, "v1.1", semver.Increment(LevelMinor).String())
			})
			t.Run("increment major", func(t *testing.T) {
				semver, err := Parse("v1-notes")
				require.NoError(t, err)
				assert.Equal(t, "v2", semver.Increment(LevelMajor).String())
			})
			t.Run("decrement patch", func(t *testing.T) {
				semver, err := Parse("v1-notes")
				require.NoError(t, err)
				assert.PanicsWithValue(t, "undefined level", func() { semver.Decrement(LevelPatch) })
			})
			t.Run("decrement minor", func(t *testing.T) {
				semver, err := Parse("v1-notes")
				require.NoError(t, err)
				assert.PanicsWithValue(t, "undefined level", func() { semver.Decrement(LevelMinor) })
			})
			t.Run("decrement major", func(t *testing.T) {
				semver, err := Parse("v1-notes")
				require.NoError(t, err)
				assert.Equal(t, "v0", semver.Decrement(LevelMajor).String())
			})
		})

		t.Run("decrement zero", func(t *testing.T) {
			semver, err := Parse("v0.0.0-notes")
			require.NoError(t, err)
			assert.PanicsWithValue(t, "zero cannot be decremented", func() { semver.Decrement(LevelPatch) }, "decrementing patch")
			assert.PanicsWithValue(t, "zero cannot be decremented", func() { semver.Decrement(LevelMinor) }, "decrementing minor")
			assert.PanicsWithValue(t, "zero cannot be decremented", func() { semver.Decrement(LevelMajor) }, "decrementing major")
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
