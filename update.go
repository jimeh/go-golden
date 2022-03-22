package golden

import (
	"os"
	"strings"
)

var truthyStrings = []string{"1", "y", "t", "yes", "on", "true"}

type UpdateFunc func() bool

// EnvUpdateFunc checks if the GOLDEN_UPDATE environment variable is set to
// one of "1", "y", "t", "yes", "on", or "true".
//
// This is also the default UpdateFunc used to determine the return value of
// Update().
func EnvUpdateFunc() bool {
	env := os.Getenv("GOLDEN_UPDATE")
	for _, v := range truthyStrings {
		if strings.ToLower(env) == v {
			return true
		}
	}

	return false
}
