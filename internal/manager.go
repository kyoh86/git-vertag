package internal

type Manager struct {
	Tagger
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
	var latest Semver
	tags, err := m.retrieveTags(fetch)
	if err != nil {
		return latest, err
	}
	for _, tag := range tags {
		ver, err := ParseSemver(tag)
		if err != nil {
			continue
		}
		latest = GreaterSemver(latest, ver)
	}

	return latest, nil
}
