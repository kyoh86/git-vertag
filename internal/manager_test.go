package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	tset := func() (*bytes.Buffer, *MockRunner, *Manager) { //nolint
		buffer := &bytes.Buffer{}
		runner := &MockRunner{echo: buffer}
		manager := &Manager{Tagger: Tagger{Remote: "test", Runner: runner}}
		return buffer, runner, manager
	}

	t.Run("create ver", func(t *testing.T) {
		buf, _, man := tset()
		assert.NoError(t, man.CreateVer(NewSemver("", 1, 2, 3), nil, ""))
		assert.Equal(t, "git tag v1.2.3\n", buf.String())
	})

	t.Run("replace ver", func(t *testing.T) {
		buf, _, man := tset()
		assert.NoError(t, man.ReplaceVer(NewSemver("", 1, 2, 3), nil, ""))
		assert.Equal(t, `git tag -d v1.2.3
git tag v1.2.3
`, buf.String())
	})

	t.Run("get ver", func(t *testing.T) {
		t.Run("without tag", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("")
			ver, err := man.GetVer(false)
			assert.NoError(t, err)
			assert.Equal(t, Semver{}, ver)
			assert.Equal(t, "v0", ver.String())
			assert.Equal(t, "git tag -l\n", buf.String())
		})
		t.Run("without version tag", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("foo\nbar\n")
			ver, err := man.GetVer(false)
			assert.NoError(t, err)
			assert.Equal(t, Semver{}, ver)
			assert.Equal(t, "v0", ver.String())
			assert.Equal(t, "git tag -l\n", buf.String())
		})
		t.Run("select newest version", func(t *testing.T) {
			buf, run, man := tset()
			run.output = strings.NewReader("1.3.0\nvar\nv2\n0.3,1\nfoo\n")
			ver, err := man.GetVer(false)
			assert.NoError(t, err)
			assert.Equal(t, NewSemver("", 2), ver)
			assert.Equal(t, "v2", ver.String())
			assert.Equal(t, "git tag -l\n", buf.String())
		})
	})

}

func TestManagerFS(t *testing.T) {
	temp := func(t *testing.T) (*Manager, func()) {
		dir, err := ioutil.TempDir("", "git-vertag-test")
		if err != nil {
			t.Logf("failed to create temp dir %v", err)
			t.Skip()
		}
		tag := NewTagger()
		tag.Workdir = dir
		return &Manager{Tagger: tag}, func() { os.RemoveAll(dir) }
	}
	init := func(t *testing.T) (*Manager, func()) {
		man, tear := temp(t)
		if err := man.Tagger.run(nil, "init"); err != nil {
			t.Logf("failed to git init %v", err)
			t.Skip()
		}
		if err := man.Tagger.run(nil, "commit", "--allow-empty", "-m", "init"); err != nil {
			t.Logf("failed to create first commit %v", err)
			t.Skip()
		}
		return man, tear
	}

	t.Run("initial", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			assert.NoError(t, man.CreateVer(NewSemver("", 0, 0, 1), nil, ""))
			ver, err := man.GetVer(false)
			assert.NoError(t, err)
			assert.Equal(t, "v0.0.1", ver.String())
		})

		t.Run("get", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			ver, err := man.GetVer(false)
			assert.NoError(t, err)
			assert.Equal(t, "v0", ver.String())
		})

		t.Run("replace", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			assert.Error(t, man.ReplaceVer(NewSemver("", 0, 0, 1), nil, ""))
		})

		t.Run("delete", func(t *testing.T) {
			man, tear := init(t)
			defer tear()
			assert.Error(t, man.DeleteVer(NewSemver("", 0, 0, 1)))
		})
	})

	t.Run("empty dir", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			assert.Error(t, man.CreateVer(NewSemver("", 1), nil, ""))
		})

		t.Run("get", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			_, err := man.GetVer(false)
			assert.Error(t, err)
		})

		t.Run("replace", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			assert.Error(t, man.ReplaceVer(NewSemver("", 1), nil, ""))
		})

		t.Run("delete", func(t *testing.T) {
			man, tear := temp(t)
			defer tear()
			assert.Error(t, man.DeleteVer(NewSemver("", 1)))
		})
	})
}
