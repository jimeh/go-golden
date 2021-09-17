package golden

import (
	"regexp"
	"strings"
)

var (
	whitespaceChars = regexp.MustCompile(`\s`)
	illegalChars    = regexp.MustCompile(`[\/\?<>\\:\*\|"]`)
	controlChars    = regexp.MustCompile(`[\x00-\x1f\x80-\x9f]`)
	reservedNames   = regexp.MustCompile(`^\.+$`)
	winReserved     = regexp.MustCompile(
		`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])(\..*)?$`,
	)
)

func sanitizeFilename(name string) string {
	if reservedNames.MatchString(name) || winReserved.MatchString(name) {
		var b []byte
		for i := 0; i < len(name); i++ {
			b = append(b, byte('_'))
		}

		return string(b)
	}

	r := strings.TrimRight(name, ". ")
	r = whitespaceChars.ReplaceAllString(r, "_")
	r = illegalChars.ReplaceAllString(r, "_")
	r = controlChars.ReplaceAllString(r, "_")

	return r
}
