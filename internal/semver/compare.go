package semver

import "strings"

// Compare returns an integer comparing two Semvers lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(v1, v2 Semver) int {
	/* SPEC:
	Precedence refers to how versions are compared to each other when ordered.
	Precedence MUST be calculated by separating the version into major, minor,
	patch and pre-release identifiers in that order (Build metadata does not
	figure into precedence).
	Precedence is determined by the first difference when comparing each of
	these identifiers from left to right as follows: Major, minor, and patch
	versions are always compared numerically.

	Example: 1.0.0 < 2.0.0 < 2.1.0 < 2.1.1.

	When major, minor, and patch are equal, a pre-release version has lower
	precedence than a normal version.

	Example: 1.0.0-alpha < 1.0.0.

	Precedence for two pre-release versions with the same major, minor, and
	patch version MUST be determined by comparing each dot separated identifier
	from left to right until a difference is found as follows: identifiers
	consisting of only digits are compared numerically and identifiers with
	letters or hyphens are compared lexically in ASCII sort order.
	Numeric identifiers always have lower precedence than non-numeric identifiers.
	A larger set of pre-release fields has a higher precedence than a smaller set,
	if all of the preceding identifiers are equal.

	Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta
		< 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
	*/

	if v1.Major < v2.Major {
		return -1
	}
	if v1.Major > v2.Major {
		return 1
	}
	if v1.Minor < v2.Minor {
		return -1
	}
	if v1.Minor > v2.Minor {
		return 1
	}
	if v1.Patch < v2.Patch {
		return -1
	}
	if v1.Patch > v2.Patch {
		return 1
	}

	if len(v1.PreRelease) == 0 {
		if len(v2.PreRelease) == 0 {
			return 0
		}
		return 1
	} else if len(v2.PreRelease) == 0 {
		return -1
	}

	for i, vp1 := range v1.PreRelease {
		if len(v2.PreRelease) <= i {
			return 1
		}
		vp2 := v2.PreRelease[i]
		if vp1.isNum {
			if vp2.isNum {
				// identifiers consisting of only digits are compared numerically
				switch {
				case vp1.num < vp2.num:
					return -1
				case vp1.num > vp2.num:
					return 1
				default:
					return 0
				}
			} else {
				// SPEC: Numeric identifiers always have lower precedence than non-numeric identifiers.
				return -1
			}
		}
		if vp2.isNum {
			// SPEC: Numeric identifiers always have lower precedence than non-numeric identifiers.
			return 1
		}
		// SPEC: letters or hyphens are compared lexically in ASCII sort order.
		lex := strings.Compare(vp1.str, vp2.str)
		if lex != 0 {
			return lex
		}
	}

	if len(v2.PreRelease) > len(v1.PreRelease) {
		return -1
	}

	return 0
}
