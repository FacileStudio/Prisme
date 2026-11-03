package spaces

import "testing"

func TestSlugifyKeepsAliasSafeCharacters(t *testing.T) {
	cases := map[string]string{
		"Facile Studio":       "facile-studio",
		"  Équipe -- Design ": "quipe-design",
		"!!!":                 "space",
		"general":             "general",
	}
	for name, want := range cases {
		if got := Slugify(name, "space"); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", name, got, want)
		}
	}
}
