package entities

import (
	"golang.org/x/mod/semver"
	"time"
)

type Version struct {
	Date            time.Time `gorm:"Type:date"`
	SemanticVersion string
}

func (v *Version) IsNewer(lastVersion *Version) bool {
	if v.Date.After(lastVersion.Date) {
		return true
	}
	compareVersion := semver.Compare(v.SemanticVersion, lastVersion.SemanticVersion)
	if compareVersion == 1 {
		return true
	}
	return false
}
