package internal

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
)

/* SPEC: https://version.org/ */

func mustPRVer(t *testing.T, s string) semver.PRVersion {
	t.Helper()
	v, err := semver.NewPRVersion(s)
	if err != nil {
		t.Error(err)
	}
	return v
}

func TextIncrementPre(t *testing.T) {
	t.Run("increment pre-release", func(t *testing.T) {
		version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
		assert.NoError(t, err)
		v, err := NewUpdater(version).
			Pre().Version()
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3-pre-release.5", v.String())
	})
	t.Run("incrementing pre-release, remaining words should be dropped", func(t *testing.T) {
		version, err := semver.Parse("1.2.3-pre-release.4.x.y")
		assert.NoError(t, err)
		v, err := NewUpdater(version).
			Pre().Version()
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3-pre-release.5", v.String())
	})
	t.Run("incrementing pre-release without numeric suffix", func(t *testing.T) {
		version, err := semver.Parse("1.2.3-pre-release")
		assert.NoError(t, err)
		v, err := NewUpdater(version).
			Pre().Version()
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3-pre-release.2", v.String())
	})
	t.Run("incrementing pre-release without pre-release never changes", func(t *testing.T) {
		version, err := semver.Parse("1.2.3")
		assert.NoError(t, err)
		v, err := NewUpdater(version).
			Pre().Version()
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3", v.String())
	})
}

func TestUpdater(t *testing.T) {
	t.Run("source: pre-release+build-ver", func(t *testing.T) {
		t.Run("set build", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).
				Build("build-ver", ("6")).Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.3-pre-release.4+build-ver.6", v.String())
		})
		t.Run("set pre-release and build", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).
				Pre(mustPRVer(t, "beta"), mustPRVer(t, "6")).
				Build(("build-ver"), ("7")).Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.3-beta.6+build-ver.7", v.String())
		})
		t.Run("set pre-release", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).
				Pre(mustPRVer(t, "beta"), mustPRVer(t, "6")).Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.3-beta.6", v.String())
		})
		t.Run("increment pre-release", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).
				Pre().Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.3-pre-release.5", v.String())
		})
		t.Run("incrementing pre-release, remaining words should be dropped", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4.x.y")
			assert.NoError(t, err)
			v, err := NewUpdater(version).
				Pre().Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.3-pre-release.5", v.String())
		})
		t.Run("increment patch", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).Patch().Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.4", v.String())
		})
		t.Run("increment minor", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).Minor().Version()
			assert.Equal(t, "1.3.0", v.String())
			assert.NoError(t, err)
		})
		t.Run("increment major", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4+build-ver.5")
			assert.NoError(t, err)
			v, err := NewUpdater(version).Major().Version()
			assert.NoError(t, err)
			assert.Equal(t, "2.0.0", v.String())
		})
	})

	t.Run("source: pre-release", func(t *testing.T) {
		t.Run("increment patch", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4")
			assert.NoError(t, err)
			v, err := NewUpdater(version).Patch().Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.2.4", v.String())
		})
		t.Run("increment minor", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4")
			assert.NoError(t, err)
			v, err := NewUpdater(version).Minor().Version()
			assert.NoError(t, err)
			assert.Equal(t, "1.3.0", v.String())
		})
		t.Run("increment major", func(t *testing.T) {
			version, err := semver.Parse("1.2.3-pre-release.4")
			assert.NoError(t, err)
			v, err := NewUpdater(version).Major().Version()
			assert.NoError(t, err)
			assert.Equal(t, "2.0.0", v.String())
		})
	})

}
