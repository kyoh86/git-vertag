package semver

import "errors"

func (v Semver) Update() Updater {
	return NewUpdater(v)
}

func NewUpdater(v Semver) Updater {
	return &implUpdater{src: v, ver: v}
}

// Updater will update version without rewinds.
type Updater interface {
	UpdatePreRelease
	Major() UpdatePreRelease
	Minor() UpdatePreRelease
	Patch() UpdatePreRelease
}

type UpdatePreRelease interface {
	UpdateBuild
	PreRelease(...PreReleaseID) UpdateBuild
}

type UpdateBuild interface {
	Applier
	Build(...BuildID) Applier
}

type Applier interface {
	Apply() (Semver, error)
}

type implUpdater struct {
	src Semver
	ver Semver
}

func (i *implUpdater) Major() UpdatePreRelease {
	i.ver.Major += 1
	i.ver.Minor = 0
	i.ver.Patch = 0
	i.ver.PreRelease = nil
	i.ver.Build = nil
	return i
}

func (i *implUpdater) Minor() UpdatePreRelease {
	i.ver.Minor += 1
	i.ver.Patch = 0
	i.ver.PreRelease = nil
	i.ver.Build = nil
	return i
}

func (i *implUpdater) Patch() UpdatePreRelease {
	i.ver.Patch += 1
	i.ver.PreRelease = nil
	i.ver.Build = nil
	return i
}

func (i *implUpdater) PreRelease(p ...PreReleaseID) UpdateBuild {
	i.ver.PreRelease = p
	i.ver.Build = nil
	return i
}

func (i *implUpdater) Build(b ...BuildID) Applier {
	i.ver.Build = b
	return i
}

func (i *implUpdater) Apply() (Semver, error) {
	if i.src.Major == i.ver.Major && i.src.Minor == i.ver.Minor && i.src.Patch == i.ver.Patch && len(i.src.PreRelease) == 0 && len(i.ver.PreRelease) > 0 {
		return Semver{}, errors.New("putting pre-release ID rewinds version order")
	}
	return i.ver, nil
}
