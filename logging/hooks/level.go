package hooks

import "github.com/sirupsen/logrus"

// LevelString converts the Level to a string.
// WarnLevel becomes "warn".
func LevelString(level logrus.Level) string {
	if level == logrus.WarnLevel {
		return "warn"
	}
	return level.String()
}
