package agents

import "testing"

func TestChooseLocalpartPrefersTheSlugThenTheOwnerSuffix(t *testing.T) {
	cases := []struct {
		taken []string
		want  string
	}{
		{nil, "pixel"},
		{[]string{"pixel"}, "pixel-yann"},
		{[]string{"pixel", "pixel-yann"}, "pixel-yann2"},
		{[]string{"pixel", "pixel-yann", "pixel-yann2"}, "pixel-yann3"},
	}
	for _, c := range cases {
		held := map[string]bool{}
		for _, localpart := range c.taken {
			held[localpart] = true
		}
		got := chooseLocalpart("pixel", "yann", func(localpart string) bool { return held[localpart] })
		if got != c.want {
			t.Errorf("taken %v: chooseLocalpart = %q, want %q", c.taken, got, c.want)
		}
	}
	if got := chooseLocalpart("pixel", "yann", func(string) bool { return true }); got != "" {
		t.Errorf("everything taken: chooseLocalpart = %q, want empty", got)
	}
}
