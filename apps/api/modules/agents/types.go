package agents

// CreateRequest is the body of POST /agents.
type CreateRequest struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// Login is the agent's Matrix login, returned by POST /agents.
type Login struct {
	Homeserver  string `json:"homeserver"`
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
}

// Agent is one of the caller's agents, as GET /agents lists them.
type Agent struct {
	UserID  string   `json:"user_id"`
	Name    string   `json:"name"`
	RoomIDs []string `json:"room_ids"`
}

// RoomsRequest is the body of PUT /agents/{user_id}/rooms.
type RoomsRequest struct {
	RoomIDs []string `json:"room_ids"`
}

// RoomsResponse lists the rooms the agent is joined to afterwards.
type RoomsResponse struct {
	RoomIDs []string `json:"room_ids"`
}
