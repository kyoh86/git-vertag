package internal

import "github.com/kyoh86/git-vertag/internal/semver"

type Manager struct {
	Tagger Tagger
}

func (m *Manager) DeleteVer(v semver.Semver) error {
	if err := m.Tagger.DeleteTag("v" + v.String()); err != nil {
		return err
	}
	return nil
}

func (m *Manager) CreateVer(v semver.Semver, message []string, file string) error {
	if err := m.Tagger.CreateTag("v"+v.String(), message, file); err != nil {
		return err
	}
	return nil
}

func (m *Manager) ReplaceVer(v semver.Semver, message []string, file string) error {
	if err := m.DeleteVer(v); err != nil {
		return err
	}
	if err := m.CreateVer(v, message, file); err != nil {
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
		ver, err := semver.ParseTolerant(tag)
		if err != nil {
			continue
		}
		if semver.Compare(latest, ver) < 0 {
			latest = ver
		}
	}

	return latest, nil
}
