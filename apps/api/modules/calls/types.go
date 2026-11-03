package calls

// Request is the body of POST /calls. Join skips the call message, for
// someone joining a call already announced in the room.
type Request struct {
	RoomID string `json:"room_id"`
	Join   bool   `json:"join,omitempty"`
}

// Response is what a LiveKit client needs to join, plus Echo's own link.
type Response struct {
	CallID     string `json:"call_id"`
	LiveKitURL string `json:"livekit_url"`
	Token      string `json:"token"`
	JoinURL    string `json:"join_url"`
}

// CallContent is the studio.facile.call block of a call message.
type CallContent struct {
	CallID    string `json:"call_id"`
	JoinURL   string `json:"join_url"`
	StartedBy string `json:"started_by"`
}
