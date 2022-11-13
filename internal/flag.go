package internal

import (
	"strings"

	"github.com/blang/semver/v4"
)

type PRVersionFlag semver.PRVersion

func (f *PRVersionFlag) Set(s string) error {
	v, err := semver.NewPRVersion(s)
	if err != nil {
		return err
	}
	*f = PRVersionFlag(v)
	return nil
}

func (f PRVersionFlag) String() string {
	return (semver.PRVersion(f)).String()
}

type PreReleaseFlag []semver.PRVersion

func (f PreReleaseFlag) IsCumulative() bool { return true }

func (f *PreReleaseFlag) Set(s string) error {

	v, err := semver.NewPRVersion(s)
	if err != nil {
		return err
	}
	*f = append(*f, v)
	return nil
}

func (f PreReleaseFlag) String() string {
	items := make([]string, 0, len(f))
	for _, v := range f {
		items = append(items, v.String())
	}
	return strings.Join(items, ",")
}

func (f PreReleaseFlag) List() []semver.PRVersion {
	return []semver.PRVersion(f)
}

type BuildVersionFlag string

func (f *BuildVersionFlag) Set(s string) error {
	v, err := semver.NewBuildVersion(s)
	if err != nil {
		return err
	}
	*f = BuildVersionFlag(v)
	return nil
}

func (f BuildVersionFlag) String() string {
	return string(f)
}

type BuildFlag []string

func (f BuildFlag) IsCumulative() bool { return true }

func (f *BuildFlag) Set(s string) error {
	v, err := semver.NewBuildVersion(s)
	if err != nil {
		return err
	}
	*f = append(*f, v)
	return nil
}

func (f BuildFlag) String() string {
	return strings.Join([]string(f), ".")
}

func (f BuildFlag) List() []string {
	return []string(f)
}
