package calls

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

const keyLabel = "prisme-call-e2ee:"

// keyMaterial turns PRISME_SECRET_KEY into the HMAC key for call keys. The
// label keeps it apart from the seal key, which hashes the bare secret.
func keyMaterial(secret string) []byte {
	sum := sha256.Sum256([]byte(keyLabel + secret))
	return sum[:]
}

// roomKey is the end-to-end encryption key of a Matrix room's calls: every
// member gets the same key, and Echo never sees it.
func roomKey(material []byte, roomID string) string {
	mac := hmac.New(sha256.New, material)
	mac.Write([]byte(keyLabel + roomID))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
