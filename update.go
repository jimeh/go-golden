package golden

import (
	"flag"
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

var (
	updateFlagSet *flag.FlagSet
	updateFlag    bool
)

// UpdateFunc returns a function that checks a -update flag is set.
func FlagUpdateFunc() bool {
	if updateFlagSet == nil {
		updateFlagSet = flag.NewFlagSet("golden", flag.ContinueOnError)
		updateFlagSet.BoolVar(&updateFlag,
			"update", false, "update golden files",
		)
	}

	_ = updateFlagSet.Parse(os.Args[1:])

	return false
}
