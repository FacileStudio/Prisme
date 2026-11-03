package calls

import "testing"

func TestRoomKeyIsStablePerRoomAndSecret(t *testing.T) {
	material := keyMaterial("a secret that is long enough to boot")
	first := roomKey(material, "!a:prisme.test")
	if first != roomKey(keyMaterial("a secret that is long enough to boot"), "!a:prisme.test") {
		t.Fatal("the same room got two keys")
	}
	if first == roomKey(material, "!b:prisme.test") {
		t.Fatal("two rooms share a key")
	}
	if first == roomKey(keyMaterial("another secret, just as long"), "!a:prisme.test") {
		t.Fatal("two secrets give the same key")
	}
	if len(first) != 43 {
		t.Fatalf("key length = %d, want 43 base64url characters", len(first))
	}
}
