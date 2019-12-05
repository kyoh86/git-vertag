package internal

import (
	"bufio"
	"bytes"
	"io"
)

type Manager struct {
	Command TagCommand
	Workdir string
	Push    bool
}

func (m *Manager) run(w io.Writer, args ...string) error {
	if m.Workdir != "" {
		return m.run(w, append([]string{"-C", m.Workdir}, args...)...)
	} else {
		return m.run(w, args...)
	}
}

func (m *Manager) removeTag(tag string) error {
	if err := m.run(nil, "tag", "-d", tag); err != nil {
		return err
	}

	if m.Push {
		// UNDONE: remote name (not only origin)
		if err := m.run(nil, "push", "origin", ":"+tag); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createTag(tag string, message []string, file string) error {
	args := []string{"tag"}
	for _, m := range message {
		args = append(args, "--message", m)
	}
	if file != "" {
		args = append(args, "--file", file)
	}

	if err := m.run(nil, append(args, tag)...); err != nil {
		return err
	}

	if m.Push {
		// UNDONE: remote name (not only origin)
		if err := m.run(nil, "push", "origin", tag); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) DeleteVer(v Semver) error {
	if err := m.removeTag(v.String()); err != nil {
		return err
	}
	if v.Level() == LevelPatch {
		_ = m.removeTag(v.MinorString())
	}
	if v.Level() != LevelMajor {
		_ = m.removeTag(v.MajorString())
	}
	return nil
}

func (m *Manager) CreateVer(v Semver, message []string, file string) error {
	if err := m.createTag(v.String(), message, file); err != nil {
		return err
	}
	if v.Level() == LevelPatch {
		if err := m.createTag(v.MinorString(), message, file); err != nil {
			return err
		}
	}
	if v.Level() != LevelMajor {
		if err := m.createTag(v.MajorString(), message, file); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) ReplaceVer(v Semver, message []string, file string) error {
	if err := m.DeleteVer(v); err != nil {
		return err
	}
	if err := m.CreateVer(v, message, file); err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetVer(fetch bool) (Semver, error) {
	latest := Semver{}
	if fetch {
		if err := m.run(nil, "fetch", "--tags"); err != nil {
			return latest, err
		}
	}
	var stdout bytes.Buffer
	if err := m.run(&stdout, "tag", "-l"); err != nil {
		return latest, err
	}

	stream := bufio.NewScanner(&stdout)
	for stream.Scan() {
		ver, err := ParseSemver(stream.Text())
		if err != nil {
			continue
		}
		latest = GreaterSemver(latest, ver)
	}

	return latest, nil
}
