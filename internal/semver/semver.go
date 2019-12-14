package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/* SPEC: https://semver.org/ */

type Semver struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	PreRelease PreRelease
	Build      Build
}

type PreRelease []string
type Build []string

func (s Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d%s%s", s.Major, s.Minor, s.Patch, s.PreRelease, s.Build)
}

var semverRegex = regexp.MustCompile(
	`^v?` +
		`(?P<major>\d+)` +
		`(?:\.(?P<minor>\d+))?` +
		`(?:\.(?P<patch>\d+))?` +
		`(?:-(?P<prerelease>[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?` +
		`(?:\+(?P<build>[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?` +
		`$`)

func ParseSemver(s string) (Semver, error) {
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
			ver.PreRelease = PreRelease(strings.Split(match[i], "."))
		case "build":
			ver.Build = Build(strings.Split(match[i], "."))
		}
	}
	return ver, nil
}

func Greater(v1, v2 Semver) Semver {
	/* SPEC:
	Precedence for two pre-release versions with the same major, minor,
	and patch version MUST be determined by comparing each dot separated
	identifier from left to right until a difference is found as follows:
	identifiers consisting of only digits are compared numerically and
	identifiers with letters or hyphens are compared lexically in ASCII
	sort order.
	Numeric identifiers always have lower precedence than
	non-numeric identifiers.
	A larger set of pre-release fields has a higher precedence than
	a smaller set, if all of the preceding identifiers are equal.
	Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta
		< 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
	*/

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
	// TODO: pre-release
	/* SPEC:
	When major, minor, and patch are equal,
	a pre-release version has lower precedence than a normal version.
	Example: 1.0.0-alpha < 1.0.0.

	Precedence for two pre-release versions with the same major, minor,
	and patch version MUST be determined by comparing each dot separated identifier
	from left to right until a difference is found as follows:
	identifiers consisting of only digits are compared numerically and identifiers
	with letters or hyphens are compared lexically in ASCII sort order.
	Numeric identifiers always have lower precedence than non-numeric identifiers.
	A larger set of pre-release fields has a higher precedence than a smaller set,
	if all of the preceding identifiers are equal.
	Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta
	  < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
	*/

	// IGNORE: build.
	/* SPEC:
	 * Build metadata MUST be ignored when determining version precedence.
	 * Thus two versions that differ only in the build metadata,
	 * have the same precedence.
	 */
	return v1
}
