package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Level uint8

const (
	LevelMajor Level = iota
	LevelMinor
	LevelPatch
)

type Semver struct {
	major uint64
	minor uint64
	patch uint64
	level Level
	notes string
}

func NewSemver(note string, levels ...uint64) (v Semver) {
	if note != "" {
		v.notes = "-" + note
	}
	for i, l := range levels {
		v.level = Level(i)
		switch i {
		case 0:
			v.major = l
		case 1:
			v.minor = l
		case 2:
			v.patch = l
		default:
			panic("invalid operation")
		}
	}
	return
}

func (s Semver) Notes() string {
	return strings.TrimPrefix(s.notes, "-")
}

func (s Semver) Level() Level {
	return s.level
}

func (s Semver) Increment(l Level) (retVer Semver) {
	if l > LevelPatch {
		panic("invalid operation")
	}
	retVer = s
	retVer.notes = ""
	if s.level < l {
		retVer.level = l
	}
	switch l {
	case LevelMajor:
		retVer.major++
		retVer.minor = 0
		retVer.patch = 0
	case LevelMinor:
		retVer.minor++
		retVer.patch = 0
	case LevelPatch:
		retVer.patch++
	}
	return
}

func (s Semver) Decrement(l Level) (retVer Semver) {
	if l > LevelPatch {
		panic("invalid operation")
	}
	if s.level < l {
		panic("undefined level")
	}
	retVer = s
	retVer.notes = ""
	switch l {
	case LevelMajor:
		if retVer.major == 0 {
			panic("zero cannot be decremented")
		}
		retVer.major--
		retVer.minor = 0
		retVer.patch = 0
	case LevelMinor:
		if retVer.minor == 0 {
			panic("zero cannot be decremented")
		}
		retVer.minor--
		retVer.patch = 0
	case LevelPatch:
		if retVer.patch == 0 {
			panic("zero cannot be decremented")
		}
		retVer.patch--
	}
	return
}

func (s Semver) Truncate(l Level) (retVer Semver) {
	retVer = s
	retVer.notes = ""
	switch l {
	case LevelMajor:
		retVer.minor = 0
		retVer.patch = 0
	case LevelMinor:
		retVer.patch = 0
	}
	retVer.level = l
	return
}

func (s Semver) Major() Semver {
	return s.Truncate(LevelMajor)
}

func (s Semver) Minor() Semver {
	return s.Truncate(LevelMinor)
}

func (s Semver) Patch() Semver {
	return s.Truncate(LevelPatch)
}

func (s Semver) MajorString() string {
	return s.Major().String()
}

func (s Semver) MinorString() string {
	return s.Minor().String()
}

func (s Semver) PatchString() string {
	return s.Patch().String()
}

// NOTE: Semver must not have "Full()" method. We can use "String()".
// Because "Full" sounds me that will convert "v1.2" to "v1.2.0".

func (s Semver) String() string {
	switch s.level {
	case LevelPatch:
		return fmt.Sprintf("v%d.%d.%d%s", s.major, s.minor, s.patch, s.notes)
	case LevelMinor:
		return fmt.Sprintf("v%d.%d%s", s.major, s.minor, s.notes)
	case LevelMajor:
		return fmt.Sprintf("v%d%s", s.major, s.notes)
	default:
		panic("invalid operatoin")
	}
}

var semverRegex = regexp.MustCompile(`^v?(?P<major>\d+)(\.(?P<minor>\d+))?(\.(?P<patch>\d+))?(?P<notes>-.*)?$`)

func Parse(s string) (Semver, error) {
	ver := Semver{}
	match := semverRegex.FindStringSubmatch(s)
	if len(match) == 0 {
		return ver, errors.New("invalid version syntax")
	}
	for i, name := range semverRegex.SubexpNames() {
		if match[i] == "" {
			continue
		}
		switch name {
		case "major":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.major = level
			ver.level = LevelMajor
		case "minor":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.minor = level
			ver.level = LevelMinor
		case "patch":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.patch = level
			ver.level = LevelPatch
		case "notes":
			ver.notes = match[i]
		}
	}
	return ver, nil
}

func GreaterSemver(v1, v2 Semver) Semver {
	if v1.major < v2.major {
		return v2
	}
	if v1.major > v2.major {
		return v1
	}
	if v1.minor < v2.minor {
		return v2
	}
	if v1.minor > v2.minor {
		return v1
	}
	if v1.patch < v2.patch {
		return v2
	}
	if v1.patch > v2.patch {
		return v1
	}
	return v1
}
