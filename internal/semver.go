package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Semver :
type Semver struct {
	Major int
	Minor int
	Patch int
}

func (s *Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}

var semverRegex = regexp.MustCompile(`^v?(?P<major>\d+)(\.(?P<minor>\d+))?(\.(?P<patch>\d+))?(?:-.*)?$`)

func ParseSemver(s string) (*Semver, error) {
	match := semverRegex.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil, errors.New("invalid version syntax")
	}
	result := map[string]int{}
	for i, name := range semverRegex.SubexpNames() {
		if i == 0 {
			continue
		}
		if i < len(match) {
			result[name], _ = strconv.Atoi(match[i])
		}
	}
	return &Semver{
		Major: result["major"],
		Minor: result["minor"],
		Patch: result["patch"],
	}, nil
}
func GreaterSemver(v1, v2 *Semver) *Semver {
	if v1 == nil {
		return v2
	}

	if v1.Major < v2.Major {
		return v2
	}
	if v1.Major > v2.Major {
		return v1
	}
	if v1.Minor < v2.Minor {
		return v2
	}
	if v1.Minor > v2.Minor {
		return v1
	}
	if v1.Patch < v2.Patch {
		return v2
	}
	if v1.Patch > v2.Patch {
		return v1
	}
	return v1
}
