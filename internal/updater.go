package internal

import (
	"errors"

	"github.com/blang/semver/v4"
)

func NewUpdater(v semver.Version) Updater {
	return &implUpdater{ver: v, pre: len(v.Pre) > 0}
}

// Updater will update version without rewinds.
type Updater interface {
	UpdatePre
	Major() UpdatePre
	Minor() UpdatePre
	Patch() UpdatePre
}

type UpdatePre interface {
	UpdateBuild
	Pre(...semver.PRVersion) UpdateBuild
}

type UpdateBuild interface {
	Version
	Build(...string) Version
}

type Version interface {
	Version() (semver.Version, error)
}

type implUpdater struct {
	ver semver.Version
	pre bool
}

func (i implUpdater) Major() UpdatePre {
	return implUpdater{
		ver: semver.Version{
			Major: i.ver.Major + 1,
		},
		pre: true,
	}
}

func (i implUpdater) Minor() UpdatePre {
	return implUpdater{
		ver: semver.Version{
			Major: i.ver.Major,
			Minor: i.ver.Minor + 1,
		},
		pre: true,
	}
}

func (i implUpdater) Patch() UpdatePre {
	next := i
	next.ver.Patch += 1
	next.ver.Pre = nil
	next.ver.Build = nil
	next.pre = true
	return next
}

func (i implUpdater) Pre(p ...semver.PRVersion) UpdateBuild {
	next := i
	next.ver.Pre = p
	next.ver.Build = nil
	return next
}

func (i implUpdater) Build(b ...string) Version {
	next := i
	next.ver.Build = b
	return next
}

func (i implUpdater) Version() (semver.Version, error) {
	if !i.pre && len(i.ver.Pre) > 0 {
		return semver.Version{}, errors.New("putting pre-release ID rewinds version order")
	}
	return i.ver, nil
}
