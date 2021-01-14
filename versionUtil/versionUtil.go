package versionUtil

import (
	timeutil "github.com/gyf841010/pz-infra-new/timeUtil"
)

// Util to Generate a New Version Code
func GenerateVersionCode() int64 {
	return timeutil.CurrentUnix()
}

// Compare returns an integer comparing two version Code
// The result will be 0 if versionA==versionB, -1 if versionA < versionB, and +1 if versionA > versionB.
func CompareVersionCode(versionA, versionB int64) int {
	if versionA < versionB {
		return -1
	}
	if versionA == versionB {
		return 0
	}
	return 1
}
