package internal

import (
	"bufio"
	"bytes"
	"io"
)

type Tagger struct {
	Runner  Runner
	Workdir string
	PushTo  string
}

func NewTagger() Tagger {
	return Tagger{
		Runner: NewGitRunner(),
	}
}

func (t *Tagger) run(w io.Writer, args ...string) error {
	if t.Workdir != "" {
		return t.Runner.Run(w, append([]string{"-C", t.Workdir}, args...)...)
	} else {
		return t.Runner.Run(w, args...)
	}
}

func (t *Tagger) CreateTag(tag string, message []string, file string) error {
	args := []string{"tag"}
	for _, t := range message {
		args = append(args, "--message", t)
	}
	if file != "" {
		args = append(args, "--file", file)
	}

	if err := t.run(nil, append(args, tag)...); err != nil {
		return err
	}

	if t.PushTo != "" {
		if err := t.run(nil, "push", t.PushTo, tag); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tagger) DeleteTag(tag string) error {
	if err := t.run(nil, "tag", "-d", tag); err != nil {
		return err
	}

	if t.PushTo != "" {
		if err := t.run(nil, "push", t.PushTo, ":"+tag); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tagger) GetTags(fetch bool) ([]string, error) {
	if fetch {
		if err := t.run(nil, "fetch", "--tags"); err != nil {
			return nil, err
		}
	}
	var buf bytes.Buffer
	if err := t.run(&buf, "tag", "-l"); err != nil {
		return nil, err
	}
	var tags []string
	stream := bufio.NewScanner(&buf)
	for stream.Scan() {
		tags = append(tags, stream.Text())
	}
	return tags, nil
}
