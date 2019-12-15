package internal

import (
	"strings"

	"github.com/kyoh86/git-vertag/internal/semver"
)

type Manager struct {
	Prefix string
	Tagger Tagger
}

func (m *Manager) DeleteVer(v semver.Semver) error {
	if err := m.Tagger.DeleteTag(m.Prefix + v.String()); err != nil {
		return err
	}
	return nil
}

func (m *Manager) CreateVer(v semver.Semver, message []string, file string) error {
	if err := m.Tagger.CreateTag(m.Prefix+v.String(), message, file); err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetVer(fetch bool) (semver.Semver, error) {
	var latest semver.Semver
	tags, err := m.Tagger.GetTags(fetch)
	if err != nil {
		return latest, err
	}
	for _, tag := range tags {
		if !strings.HasPrefix(tag, m.Prefix) {
			continue
		}
		tag = strings.TrimPrefix(tag, m.Prefix)
		ver, err := semver.Parse(tag)
		if err != nil {
			continue
		}
		if semver.Compare(latest, ver) < 0 {
			latest = ver
		}
	}

	return latest, nil
}
