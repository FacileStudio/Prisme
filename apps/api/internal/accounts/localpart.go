package accounts

import (
	"strconv"
	"strings"
)

const maxLocalpart = 40

// Localpart derives a Matrix localpart from an email address: the part before
// the @, lower-cased, with every character Matrix ids refuse turned into a
// dot. It never returns something in the agent namespace.
func Localpart(email string) string {
	local, _, _ := strings.Cut(strings.ToLower(strings.TrimSpace(email)), "@")
	var b strings.Builder
	for _, r := range local {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
			continue
		}
		b.WriteRune('.')
	}
	out := strings.Trim(b.String(), "._-")
	if len(out) > maxLocalpart {
		out = strings.Trim(out[:maxLocalpart], "._-")
	}
	if out == "" {
		return "user"
	}
	if strings.HasPrefix(out, "agent-") || isReserved(out) {
		return "u." + out
	}
	return out
}

// Candidate is the n-th choice for a base localpart: the base, then base2,
// base3 and so on.
func Candidate(base string, n int) string {
	if n <= 1 {
		return base
	}
	return base + strconv.Itoa(n)
}

// isReserved reports the localparts never handed to a human: the appservice
// sender, the agent namespace, and names that read as staff.
func isReserved(localpart string) bool {
	switch localpart {
	case "prisme", "admin", "root", "system", "agent":
		return true
	}
	return false
}
