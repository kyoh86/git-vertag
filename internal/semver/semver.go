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

// IsCumlative implements kingpin.repeatableFlag
func (p *PreRelease) IsCumulative() bool {
	return true
}

func (p *PreRelease) Set(s string) error {
	newp, err := ParsePreReleaseID(s)
	if err != nil {
		return err
	}
	*p = append(*p, newp)
	return nil
}

func (p PreRelease) String() string {
	if len(p) == 0 {
		return ""
	}
	b := make([]byte, 0, len(p)*2-1)
	b = append(b, p[0].str...)
	for _, pre := range p[1:] {
		b = append(b, '.')
		b = append(b, pre.String()...)
	}
	return string(b)
}

type PreReleaseID struct {
	str   string
	num   uint64
	isNum bool
}

func (p PreReleaseID) String() string {
	return p.str
}

func (p *PreReleaseID) Set(s string) error {
	newp, err := ParsePreReleaseID(s)
	if err != nil {
		return err
	}
	*p = newp
	return nil
}

type Build []BuildID

// IsCumlative implements kingpin.repeatableFlag
func (b *Build) IsCumulative() bool {
	return true
}

func (b *Build) Set(s string) error {
	newb, err := ParseBuildID(s)
	if err != nil {
		return err
	}
	*b = append(*b, newb)
	return nil
}

func (b Build) String() string {
	if len(b) == 0 {
		return ""
	}
	a := make([]byte, 0, len(b)*2-1)
	a = append(a, b[0]...)
	for _, pre := range b[1:] {
		a = append(a, '.')
		a = append(a, pre.String()...)
	}
	return string(a)
}

type BuildID string

func (b BuildID) String() string {
	return string(b)
}

func (b *BuildID) Set(s string) error {
	newp, err := ParseBuildID(s)
	if err != nil {
		return err
	}
	*b = newp
	return nil
}
