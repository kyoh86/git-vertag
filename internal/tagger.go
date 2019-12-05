package internal

import (
	"bufio"
	"bytes"
	"io"
)

type Tagger struct {
	Runner  Runner
	Workdir string
	Push    bool
}

func (t *Tagger) run(w io.Writer, args ...string) error {
	if t.Workdir != "" {
		return t.Runner.Run(w, append([]string{"-C", t.Workdir}, args...)...)
	} else {
		return t.Runner.Run(w, args...)
	}
}

func (t *Tagger) removeTag(tag string) error {
	if err := t.run(nil, "tag", "-d", tag); err != nil {
		return err
	}

	if t.Push {
		// UNDONE: remote name (not only origin)
		if err := t.run(nil, "push", "origin", ":"+tag); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tagger) createTag(tag string, message []string, file string) error {
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

	if t.Push {
		// UNDONE: remote name (not only origin)
		if err := t.run(nil, "push", "origin", tag); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tagger) retrieveTags(fetch bool) ([]string, error) {
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
