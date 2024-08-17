package internal

import (
	"strconv"
	"strings"

	"github.com/blang/semver/v4"
)

type Manager struct {
	Prefix    string
	Tagger    Tagger
	Fetch     bool
	Ancestors bool
}

func (m *Manager) ancestors(v semver.Version) []string {
	if !m.Ancestors {
		return nil
	}
	var ancs []string
	b := make([]byte, 0, 3)
	b = strconv.AppendUint(b, v.Major, 10)
	ancs = append(ancs, m.Prefix+string(b))

	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)
	ancs = append(ancs, m.Prefix+string(b))
	return ancs
}

func (m *Manager) deleteVer(v semver.Version) error {
	if err := m.Tagger.DeleteTag(m.Prefix + v.String()); err != nil {
		return err
	}
	return nil
}

func (m *Manager) createVer(v semver.Version, msg []string, file string) error {
	if err := m.Tagger.CreateTag(m.Prefix+v.String(), msg, file); err != nil {
		return err
	}
	return nil
}

func (m *Manager) getVer() (semver.Version, error) {
	var latest semver.Version
	tags, err := m.Tagger.GetTags(m.Fetch)
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
		if latest.Compare(ver) < 0 {
			latest = ver
		}
	}

	return latest, nil
}

func (m *Manager) GetVer() (string, error) {
	v, err := m.getVer()
	if err != nil {
		return "", err
	}
	return m.Prefix + v.String(), nil
}

func (m *Manager) UpdateMajor(pre []semver.PRVersion, build, msg []string, file string) (string, string, error) {
	return m.update(pre, build, msg, file, func(u Updater) UpdatePre { return u.Major() })
}

func (m *Manager) UpdateMinor(pre []semver.PRVersion, build, msg []string, file string) (string, string, error) {
	return m.update(pre, build, msg, file, func(u Updater) UpdatePre { return u.Minor() })
}

func (m *Manager) UpdatePatch(pre []semver.PRVersion, build, msg []string, file string) (string, string, error) {
	return m.update(pre, build, msg, file, func(u Updater) UpdatePre { return u.Patch() })
}

func (m *Manager) UpdatePre(pre []semver.PRVersion, build, msg []string, file string) (string, string, error) {
	return m.update(pre, build, msg, file, func(u Updater) UpdatePre { return u })
}

func (m *Manager) update(
	pre []semver.PRVersion,
	build,
	msg []string,
	file string,
	upd func(Updater) UpdatePre,
) (string, string, error) {
	cur, err := m.getVer()
	if err != nil {
		return "", "", err
	}
	next, err := upd(NewUpdater(cur)).Pre(pre...).Build(build...).Version()
	if err != nil {
		return "", "", err
	}
	if err := m.createVer(next, msg, file); err != nil {
		return "", "", err
	}
	for _, anc := range m.ancestors(next) {
		if err := m.Tagger.DeleteTag(anc); err != nil {
			// noop
			_ = err
		}
		if err := m.Tagger.CreateTag(anc, msg, file); err != nil {
			return "", "", err
		}
	}
	return m.Prefix + cur.String(), m.Prefix + next.String(), nil
}

func (m *Manager) Release(build, msg []string, file string) (string, string, error) {
	return m.release(build, msg, file, func(u Updater) UpdateBuild { return u.Pre() })
}

func (m *Manager) Build(build, msg []string, file string) (string, string, error) {
	return m.release(build, msg, file, func(u Updater) UpdateBuild { return u })
}

func (m *Manager) release(build, msg []string, file string, upd func(Updater) UpdateBuild) (string, string, error) {
	cur, err := m.getVer()
	if err != nil {
		return "", "", err
	}
	next, err := upd(NewUpdater(cur)).Build(build...).Version()
	if err != nil {
		return "", "", err
	}
	if err = m.createVer(next, msg, file); err != nil {
		return "", "", err
	}
	return m.Prefix + cur.String(), m.Prefix + next.String(), nil
}
