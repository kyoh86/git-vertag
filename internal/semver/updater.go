package semver

func (s Semver) Update() Updater {
	return NewUpdater(s)
}
func NewUpdater(s Semver) Updater {
	return nil //TODO: implement
}
type Updater interface {
	UpdatePreRelease
	Increment(level Level) UpdatePreRelease
}

type Level int

const (
	LevelMajor Level = iota
	LevelMinor
	LevelPatch
)

type UpdatePreRelease interface {
	UpdateBuild
	PreRelease(PreRelease) UpdateBuild
}

type UpdateBuild interface {
	Applier
	Build(Build) Applier
}

type Applier interface {
	Apply() Semver
}
