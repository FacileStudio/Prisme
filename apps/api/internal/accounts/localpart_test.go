package accounts

import "testing"

func TestLocalpartKeepsMatrixSafeCharactersAndAvoidsTheAgentNamespace(t *testing.T) {
	cases := map[string]string{
		"Yann.Doe@facile.studio":  "yann.doe",
		"élodie+chat@example.com": "lodie.chat",
		"agent-x@example.com":     "u.agent-x",
		"Prisme@example.com":      "u.prisme",
		"@example.com":            "user",
		"__@example.com":          "user",
		"a_b-c@example.com":       "a_b-c",
	}
	for email, want := range cases {
		if got := Localpart(email); got != want {
			t.Errorf("Localpart(%q) = %q, want %q", email, got, want)
		}
	}
	if Candidate("tom", 1) != "tom" || Candidate("tom", 3) != "tom3" {
		t.Error("Candidate does not number from the second choice")
	}
}
