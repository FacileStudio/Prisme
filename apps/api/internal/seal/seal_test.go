package seal

import "testing"

func TestSealRoundTripsAndRefusesAnotherKey(t *testing.T) {
	box, err := New("a secret that is long enough")
	if err != nil {
		t.Fatal(err)
	}
	sealed := box.Seal("syt_token")
	if sealed == "syt_token" {
		t.Fatal("the value was not encrypted")
	}
	opened, err := box.Open(sealed)
	if err != nil || opened != "syt_token" {
		t.Fatalf("Open = %q, %v", opened, err)
	}
	other, _ := New("another secret, also long")
	if _, err := other.Open(sealed); err == nil {
		t.Fatal("a different key opened the value")
	}
}
