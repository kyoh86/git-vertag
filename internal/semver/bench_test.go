package semver

import (
	"testing"
)

type formatTest struct {
	v      Semver
	result string
}

var formatTests = []formatTest{
	{Semver{1, 2, 3, nil, nil}, "1.2.3"},
	{Semver{0, 0, 1, nil, nil}, "0.0.1"},
	{Semver{0, 0, 1, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{str: "preview"}}, Build{"123", "456"}}, "0.0.1-alpha.preview+123.456"},
	{Semver{1, 2, 3, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{isNum: true, num: 1}}, Build{"123", "456"}}, "1.2.3-alpha.1+123.456"},
	{Semver{1, 2, 3, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{isNum: true, num: 1}}, nil}, "1.2.3-alpha.1"},
	{Semver{1, 2, 3, nil, Build{"123", "456"}}, "1.2.3+123.456"},
	{Semver{1, 2, 3, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{str: "b-eta"}}, Build{"123", "b-uild"}}, "1.2.3-alpha.b-eta+123.b-uild"},
	{Semver{1, 2, 3, nil, Build{"123", "b-uild"}}, "1.2.3+123.b-uild"},
	{Semver{1, 2, 3, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{str: "b-eta"}}, nil}, "1.2.3-alpha.b-eta"},
}

type compareTest struct {
	v1     Semver
	v2     Semver
	result int
}

var compareTests = []compareTest{
	{Semver{1, 0, 0, nil, nil}, Semver{1, 0, 0, nil, nil}, 0},
	{Semver{2, 0, 0, nil, nil}, Semver{1, 0, 0, nil, nil}, 1},
	{Semver{0, 1, 0, nil, nil}, Semver{0, 1, 0, nil, nil}, 0},
	{Semver{0, 2, 0, nil, nil}, Semver{0, 1, 0, nil, nil}, 1},
	{Semver{0, 0, 1, nil, nil}, Semver{0, 0, 1, nil, nil}, 0},
	{Semver{0, 0, 2, nil, nil}, Semver{0, 0, 1, nil, nil}, 1},
	{Semver{1, 2, 3, nil, nil}, Semver{1, 2, 3, nil, nil}, 0},
	{Semver{2, 2, 4, nil, nil}, Semver{1, 2, 4, nil, nil}, 1},
	{Semver{1, 3, 3, nil, nil}, Semver{1, 2, 3, nil, nil}, 1},
	{Semver{1, 2, 4, nil, nil}, Semver{1, 2, 3, nil, nil}, 1},

	// Spec Examples #11
	{Semver{1, 0, 0, nil, nil}, Semver{2, 0, 0, nil, nil}, -1},
	{Semver{2, 0, 0, nil, nil}, Semver{2, 1, 0, nil, nil}, -1},
	{Semver{2, 1, 0, nil, nil}, Semver{2, 1, 1, nil, nil}, -1},

	// Spec Examples #9
	{Semver{1, 0, 0, nil, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "alpha"}}, nil}, 1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "alpha"}}, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{isNum: true, num: 1}}, nil}, -1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{isNum: true, num: 1}}, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{str: "beta"}}, nil}, -1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "alpha"}, PreReleaseID{str: "beta"}}, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "beta"}}, nil}, -1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "beta"}}, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "beta"}, PreReleaseID{isNum: true, num: 2}}, nil}, -1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "beta"}, PreReleaseID{isNum: true, num: 2}}, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "beta"}, PreReleaseID{isNum: true, num: 11}}, nil}, -1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "beta"}, PreReleaseID{isNum: true, num: 11}}, nil}, Semver{1, 0, 0, PreRelease{PreReleaseID{str: "rc"}, PreReleaseID{isNum: true, num: 1}}, nil}, -1},
	{Semver{1, 0, 0, PreRelease{PreReleaseID{str: "rc"}, PreReleaseID{isNum: true, num: 1}}, nil}, Semver{1, 0, 0, nil, nil}, -1},

	// Ignore Build metadata
	{Semver{1, 0, 0, nil, Build{"1", "2", "3"}}, Semver{1, 0, 0, nil, nil}, 0},
}

func BenchmarkParseSimple(b *testing.B) {
	const VERSION = "0.0.1"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(VERSION)
	}
}

func BenchmarkParseComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(VERSION)
	}
}

func BenchmarkParseAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse(formatTests[n%l].result)
	}
}

func BenchmarkStringSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringLarger(b *testing.B) {
	const VERSION = "11.15.2012"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = formatTests[n%l].v.String()
	}
}

func BenchmarkCompareSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Compare(v, v)
	}
}

func BenchmarkCompareComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Compare(v, v)
	}
}

func BenchmarkCompareAverage(b *testing.B) {
	l := len(compareTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Compare(compareTests[n%l].v1, (compareTests[n%l].v2))
	}
}
