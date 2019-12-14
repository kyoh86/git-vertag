package semver

import (
	"strconv"
)

/* SPEC: https://semver.org/ */

type Semver struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	PreRelease PreRelease
	Build      Build
}

func (v Semver) String() string {
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, v.Major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Patch, 10)

	if len(v.PreRelease) > 0 {
		b = append(b, '-')
		b = append(b, v.PreRelease[0].String()...)

		for _, pre := range v.PreRelease[1:] {
			b = append(b, '.')
			b = append(b, pre.String()...)
		}
	}

	if len(v.Build) > 0 {
		b = append(b, '+')
		b = append(b, v.Build[0]...)

		for _, build := range v.Build[1:] {
			b = append(b, '.')
			b = append(b, build...)
		}
	}

	return string(b)
}

type PreRelease []PreReleaseID

type PreReleaseID struct {
	str   string
	num   uint64
	isNum bool
}

func (p PreReleaseID) String() string {
	return p.str
}

type Build []BuildID

type BuildID string

func (b BuildID) String() string {
	return string(b)
}
