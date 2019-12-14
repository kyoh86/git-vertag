package semver

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var semverRegex = regexp.MustCompile(
	`^` +
		`(?P<major>\d+)` +
		`\.(?P<minor>\d+)` +
		`\.(?P<patch>\d+)` +
		`(?:-(?P<prerelease>[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?` +
		`(?:\+(?P<build>[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?` +
		`$`)

func MustParse(s string) Semver {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

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
			ver.Major = level
		case "minor":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.Minor = level
		case "patch":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.Patch = level
		case "prerelease":
			for _, s := range strings.Split(match[i], ".") {
				ver.PreRelease = append(ver.PreRelease, parsePreReleaseID(s))
			}
		case "build":
			for _, s := range strings.Split(match[i], ".") {
				ver.Build = append(ver.Build, parseBuildID(s))
			}
		}
	}
	return ver, nil
}

var semverTolerantRegex = regexp.MustCompile(
	`^v?` +
		`(?P<major>\d+)` +
		`(?:\.(?P<minor>\d+))?` +
		`(?:\.(?P<patch>\d+))?` +
		`(?:-(?P<prerelease>[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?` +
		`(?:\+(?P<build>[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?` +
		`$`)

func MustParseTolerant(s string) Semver {
	v, err := ParseTolerant(s)
	if err != nil {
		panic(err)
	}
	return v
}

func ParseTolerant(s string) (Semver, error) {
	ver := Semver{}
	match := semverTolerantRegex.FindStringSubmatch(s)
	if len(match) == 0 {
		return ver, errors.New("invalid version syntax")
	}
	for i, name := range semverTolerantRegex.SubexpNames() {
		if match[i] == "" {
			continue
		}
		switch name {
		case "major":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.Major = level
		case "minor":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.Minor = level
		case "patch":
			level, _ := strconv.ParseUint(match[i], 10, 64)
			ver.Patch = level
		case "prerelease":
			for _, s := range strings.Split(match[i], ".") {
				ver.PreRelease = append(ver.PreRelease, parsePreReleaseID(s))
			}
		case "build":
			for _, s := range strings.Split(match[i], ".") {
				ver.Build = append(ver.Build, parseBuildID(s))
			}
		}
	}
	return ver, nil
}

// parsePreReleaseID parses a string as valid prerelease IDs
// it depends on Parse / ParseTolerant, with regexp
func parsePreReleaseID(s string) PreReleaseID {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return PreReleaseID{
			str: s,
		}
	}
	return PreReleaseID{
		str:   s,
		num:   num,
		isNum: true,
	}
}

func MustParsePreReleaseID(s string) PreReleaseID {
	p, err := ParsePreReleaseID(s)
	if err != nil {
		panic(err)
	}
	return p
}

func ParsePreReleaseID(s string) (PreReleaseID, error) {
	p := parsePreReleaseID(s)
	if p.isNum {
		return p, nil
	}
	if strings.IndexFunc(s, func(r rune) bool {
		return !('a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '-')
	}) != -1 {
		return PreReleaseID{}, errors.New("invalid format")
	}
	return p, nil
}

// parseBuildID parses a string as valid prerelease IDs
// it depends on Parse / ParseTolerant, with regexp
func parseBuildID(s string) BuildID {
	return BuildID(s)
}

func MustParseBuildID(s string) BuildID {
	p, err := ParseBuildID(s)
	if err != nil {
		panic(err)
	}
	return p
}

func ParseBuildID(s string) (BuildID, error) {
	if strings.IndexFunc(s, func(r rune) bool {
		return !('a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '-')
	}) != -1 {
		return BuildID(""), errors.New("invalid format")
	}
	return BuildID(s), nil
}
