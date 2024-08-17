package internal

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	tset := func() (*bytes.Buffer, *MockRunner, *Manager) { // nolint
		buffer := &bytes.Buffer{}
		runner := &MockRunner{echo: buffer}
		manager := &Manager{Prefix: "test", Tagger: Tagger{Runner: runner}}
		return buffer, runner, manager
	}

	t.Run("create ver", func(t *testing.T) {
		buf, _, man := tset()
		assert.NoError(t, man.createVer(semver.Version{Major: 1, Minor: 2, Patch: 3}, nil, ""))
		assert.Equal(t, "git tag test1.2.3\n", buf.String())
	})

	t.Run("get ver", func(t *testing.T) {
		t.Run("without tag", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("")
			ver, err := man.GetVer()
			assert.NoError(t, err)
			assert.Equal(t, "test0.0.0", ver)
			assert.Equal(t, "git tag -l\n", buf.String())
		})
		t.Run("without version tag", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("foo\nbar\n")
			ver, err := man.GetVer()
			assert.NoError(t, err)
			assert.Equal(t, "test0.0.0", ver)
			assert.Equal(t, "git tag -l\n", buf.String())
		})
		t.Run("select newest version", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.3.0\nvar\ntest2.0.0\ntest0.3,1\nfoo\n")
			ver, err := man.GetVer()
			assert.NoError(t, err)
			assert.Equal(t, "test2.0.0", ver)
			assert.Equal(t, "git tag -l\n", buf.String())
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("build", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.2.3-pre-release.4+build-ver.5\n")
			cur, next, err := man.Build(
				[]string{"test-bld", "2"},
				[]string{"test-msg"},
				"")
			assert.NoError(t, err)
			assert.Equal(t, "test1.2.3-pre-release.4+build-ver.5", cur)
			assert.Equal(t, "test1.2.3-pre-release.4+test-bld.2", next)
			assert.Equal(t, "git tag -l\ngit tag --message test-msg test1.2.3-pre-release.4+test-bld.2\n", buf.String())
		})
		t.Run("release", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.2.3-pre-release.4+build-ver.5\n")
			cur, next, err := man.Release(
				[]string{"test-bld", "2"},
				[]string{"test-msg"},
				"")
			assert.NoError(t, err)
			assert.Equal(t, "test1.2.3-pre-release.4+build-ver.5", cur)
			assert.Equal(t, "test1.2.3+test-bld.2", next)
			assert.Equal(t, "git tag -l\ngit tag --message test-msg test1.2.3+test-bld.2\n", buf.String())
		})
		t.Run("set pre-release", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.2.3-pre-release.4+build-ver.5\n")
			cur, next, err := man.UpdatePre(
				[]semver.PRVersion{mustPRVer(t, "test-pre"), mustPRVer(t, "1")},
				[]string{"test-bld", "2"},
				[]string{"test-msg"},
				"")
			assert.NoError(t, err)
			assert.Equal(t, "test1.2.3-pre-release.4+build-ver.5", cur)
			assert.Equal(t, "test1.2.3-test-pre.1+test-bld.2", next)
			assert.Equal(t, "git tag -l\ngit tag --message test-msg test1.2.3-test-pre.1+test-bld.2\n", buf.String())
		})
		t.Run("increment patch", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.2.3-pre-release.4+build-ver.5\n")
			cur, next, err := man.UpdatePatch(
				[]semver.PRVersion{mustPRVer(t, "test-pre"), mustPRVer(t, "1")},
				[]string{"test-bld", "2"},
				[]string{"test-msg"},
				"")
			assert.NoError(t, err)
			assert.Equal(t, "test1.2.3-pre-release.4+build-ver.5", cur)
			assert.Equal(t, "test1.2.4-test-pre.1+test-bld.2", next)
			assert.Equal(t, "git tag -l\ngit tag --message test-msg test1.2.4-test-pre.1+test-bld.2\n", buf.String())
		})
		t.Run("increment minor", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.2.3-pre-release.4+build-ver.5\n")
			cur, next, err := man.UpdateMinor(
				[]semver.PRVersion{mustPRVer(t, "test-pre"), mustPRVer(t, "1")},
				[]string{"test-bld", "2"},
				[]string{"test-msg"},
				"")
			assert.NoError(t, err)
			assert.Equal(t, "test1.2.3-pre-release.4+build-ver.5", cur)
			assert.Equal(t, "test1.3.0-test-pre.1+test-bld.2", next)
			assert.Equal(t, "git tag -l\ngit tag --message test-msg test1.3.0-test-pre.1+test-bld.2\n", buf.String())
		})
		t.Run("increment major", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("test1.2.3-pre-release.4+build-ver.5\n")
			cur, next, err := man.UpdateMajor(
				[]semver.PRVersion{mustPRVer(t, "test-pre"), mustPRVer(t, "1")},
				[]string{"test-bld", "2"},
				[]string{"test-msg"},
				"")
			assert.NoError(t, err)
			assert.Equal(t, "test1.2.3-pre-release.4+build-ver.5", cur)
			assert.Equal(t, "test2.0.0-test-pre.1+test-bld.2", next)
			assert.Equal(t, "git tag -l\ngit tag --message test-msg test2.0.0-test-pre.1+test-bld.2\n", buf.String())
		})
	})

}

func TestManagerFS(t *testing.T) {
	temp := func(t *testing.T) (*Manager, func()) {
		t.Helper()
		dir, err := os.MkdirTemp("", "git-vertag-test")
		if err != nil {
			t.Logf("failed to create temp dir %v", err)
			t.Skip()
		}
		tag := Tagger{
			Runner:  NewGitRunner(),
			Workdir: dir,
		}
		return &Manager{Tagger: tag}, func() { os.RemoveAll(dir) }
	}
	init := func(t *testing.T) (*Manager, func()) {
		t.Helper()
		man, tear := temp(t)
		if err := man.Tagger.run(true, nil, "init"); err != nil {
			t.Logf("failed to git init %v", err)
			t.Skip()
		}
		if err := man.Tagger.run(true, nil, "commit", "--allow-empty", "-m", "init"); err != nil {
			t.Logf("failed to create first commit %v", err)
			t.Skip()
		}
		return man, tear
	}

	t.Run("initial", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			assert.NoError(t, man.createVer(semver.Version{Patch: 1}, nil, ""))
			ver, err := man.GetVer()
			assert.NoError(t, err)
			assert.Equal(t, "0.0.1", ver)
		})

		t.Run("get", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			ver, err := man.GetVer()
			assert.NoError(t, err)
			assert.Equal(t, "0.0.0", ver)
		})

		t.Run("delete", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			assert.Error(t, man.deleteVer(semver.Version{Patch: 1}))
		})
	})

	t.Run("empty dir", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			assert.Error(t, man.createVer(semver.Version{Major: 1}, nil, ""))
		})

		t.Run("get", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			_, err := man.GetVer()
			assert.Error(t, err)
		})

		t.Run("delete", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			assert.Error(t, man.deleteVer(semver.Version{Major: 1}))
		})
	})
}
