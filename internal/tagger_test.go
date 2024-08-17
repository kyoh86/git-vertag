package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagger(t *testing.T) {
	tset := func() (*bytes.Buffer, *MockRunner, Tagger) {
		buffer := &bytes.Buffer{}
		runner := &MockRunner{echo: buffer}
		tagger := Tagger{Runner: runner}
		return buffer, runner, tagger
	}
	t.Run("create tag", func(t *testing.T) {
		t.Run("plain", func(t *testing.T) {
			buf, _, tag := tset()
			assert.NoError(t, tag.CreateTag("dummy", nil, ""))
			assert.Equal(t, "git tag dummy\n", buf.String())
		})

		t.Run("message text", func(t *testing.T) {
			buf, _, tag := tset()
			assert.NoError(t, tag.CreateTag("dummy", []string{"foo", "bar"}, ""))
			assert.Equal(t, "git tag --message foo --message bar dummy\n", buf.String())
		})

		t.Run("message file", func(t *testing.T) {
			buf, _, tag := tset()
			assert.NoError(t, tag.CreateTag("dummy", nil, "message.txt"))
			assert.Equal(t, "git tag --file message.txt dummy\n", buf.String())
		})

		t.Run("push", func(t *testing.T) {
			buf, _, tag := tset()
			tag.PushTo = "test"
			assert.NoError(t, tag.CreateTag("dummy", nil, ""))
			assert.Equal(t, "git tag dummy\ngit push test dummy\n", buf.String())
		})

		t.Run("workdir", func(t *testing.T) {
			buf, _, tag := tset()
			tag.PushTo = "test"
			tag.Workdir = "dir"
			assert.NoError(t, tag.CreateTag("dummy", nil, ""))
			assert.Equal(t, "git -C dir tag dummy\ngit -C dir push test dummy\n", buf.String())
		})

	})

	t.Run("delete tag", func(t *testing.T) {
		t.Run("plain", func(t *testing.T) {
			buf, _, tag := tset()
			assert.NoError(t, tag.DeleteTag("dummy"))
			assert.Equal(t, "git tag -d dummy\n", buf.String())
		})

		t.Run("push", func(t *testing.T) {
			buf, _, tag := tset()
			tag.PushTo = "test"
			assert.NoError(t, tag.DeleteTag("dummy"))
			assert.Equal(t, "git tag -d dummy\ngit push test :dummy\n", buf.String())
		})

		t.Run("workdir", func(t *testing.T) {
			buf, _, tag := tset()
			tag.PushTo = "test"
			tag.Workdir = "dir"
			assert.NoError(t, tag.DeleteTag("dummy"))
			assert.Equal(t, "git -C dir tag -d dummy\ngit -C dir push test :dummy\n", buf.String())
		})

	})

	t.Run("get tag", func(t *testing.T) {
		t.Run("plain", func(t *testing.T) {
			buf, run, tag := tset()
			run.output = strings.NewReader("foo\nbar\n")
			tags, err := tag.GetTags(false)
			assert.NoError(t, err)
			assert.Equal(t, "git tag -l\n", buf.String())
			assert.Equal(t, []string{"foo", "bar"}, tags)
		})

		t.Run("fetch", func(t *testing.T) {
			buf, run, tag := tset()
			run.output = strings.NewReader("foo\nbar\n")
			tags, err := tag.GetTags(true)
			assert.NoError(t, err)
			assert.Equal(t, "git fetch --tags\ngit tag -l\n", buf.String())
			assert.Equal(t, []string{"foo", "bar"}, tags)
		})

		t.Run("fetch in workdir", func(t *testing.T) {
			buf, run, tag := tset()
			run.output = strings.NewReader("foo\nbar\n")
			tag.Workdir = "dir"
			tags, err := tag.GetTags(true)
			assert.NoError(t, err)
			assert.Equal(t, "git -C dir fetch --tags\ngit -C dir tag -l\n", buf.String())
			assert.Equal(t, []string{"foo", "bar"}, tags)
		})

	})
}
