package version

import "github.com/Masterminds/semver"

type Version struct {
	Version string
}

func NewVersion(version string) Version {
	return Version{Version: version}
}

func (that Version) IsAtLeast(other string) (bool, error) {
	res, err := compare(that.Version, other)
	return res >= 0, err
}

func (that Version) IsLessThan(other string) (bool, error) {
	res, err := compare(that.Version, other)
	return res < 0, err
}

// If ver1 == ver2 returns 0
// If ver1 > ver2 returns 1
// If ver1 < ver2 returns -1
func compare(ver1, ver2 string) (int, error) {
	semver1, err := semver.NewVersion(ver1)
	if err != nil {
		return 0, err
	}
	semver2, err := semver.NewVersion(ver2)
	if err != nil {
		return 0, err
	}
	return semver1.Compare(semver2), nil
}
