package spaces

import (
	"strings"
)

const maxSlug = 32

// Slugify turns a display name into a lower-case, dash-separated slug that is
// safe inside a Matrix alias. It returns fallback when nothing usable is left.
func Slugify(name, fallback string) string {
	words := strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9')
	})
	out := strings.Join(words, "-")
	if len(out) > maxSlug {
		out = strings.Trim(out[:maxSlug], "-")
	}
	if out == "" {
		return fallback
	}
	return out
}
