package calls

// Request is the body of POST /calls. Join skips the call message, for
// someone joining a call already announced in the room.
type Request struct {
	RoomID string `json:"room_id"`
	Join   bool   `json:"join,omitempty"`
}

// Response is what a LiveKit client needs to join, plus Echo's own link.
// E2EEKey is the room's end-to-end encryption key: base64url of
// HMAC-SHA256 over the room id, the same for every member of the room.
type Response struct {
	CallID     string `json:"call_id"`
	LiveKitURL string `json:"livekit_url"`
	Token      string `json:"token"`
	JoinURL    string `json:"join_url"`
	E2EEKey    string `json:"e2ee_key"`
}

// CallContent is the studio.facile.call block of a call message. EndedAt
// (Unix milliseconds) and DurationMS are set once the call is over.
type CallContent struct {
	CallID     string `json:"call_id"`
	JoinURL    string `json:"join_url,omitempty"`
	StartedBy  string `json:"started_by"`
	EndedAt    int64  `json:"ended_at,omitempty"`
	DurationMS int64  `json:"duration_ms,omitempty"`
}
