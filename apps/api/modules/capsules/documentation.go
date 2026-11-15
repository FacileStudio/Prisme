package capsules

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the capsule relay for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name: "capsules",
		Description: "Messages sealed as Capsules that burn on their first read. The browser encrypts and keeps the key; " +
			"Prisme relays only ciphertext to Capsule.",
		Routes: []documentation.Route{
			{
				Method: "POST", Path: "/capsules", Summary: "Seal a capsule",
				Description: "Stores the ciphertext in Capsule as a burn-after-read capsule. url is the Capsule link without its key.",
				Auth:        "porte session", RequestBody: SealRequest{}, ResponseBody: SealResponse{},
				Errors: relayErrors(),
			},
			{
				Method: "GET", Path: "/capsules/{id}", Summary: "Check a capsule",
				Description: "exists is false once the capsule was opened or has expired.",
				Auth:        "porte session", ResponseBody: MetaResponse{},
				Errors: relayErrors(),
			},
			{
				Method: "POST", Path: "/capsules/{id}/open", Summary: "Open a capsule",
				Description: "Returns the ciphertext and burns the capsule, for everyone.",
				Auth:        "porte session", ResponseBody: OpenResponse{},
				Errors: append(relayErrors(), documentation.Error{Status: 404, Code: "not_found", Description: "Already opened or expired."}),
			},
		},
	}
}

func relayErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 400, Code: "invalid_argument", Description: "Missing ciphertext, unknown expiry or malformed id."},
		{Status: 503, Code: "unavailable", Description: "Capsule did not answer."},
	}
}
