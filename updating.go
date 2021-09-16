package golden

import "os"

var truthyStrings = []string{"1", "y", "t", "yes", "on", "true"}

type UpdatingFunc func() bool

// EnvVarUpdateFunc checks if the GOLDEN_UPDATE environment variable is set to
// one of "1", "y", "t", "yes", "on", or "true".
func EnvVarUpdatingFunc() bool {
	env := os.Getenv("GOLDEN_UPDATE")
	for _, v := range truthyStrings {
		if env == v {
			return true
		}
	}

	return false
}
