package spaces

// Channel is a room of a space.
type Channel struct {
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
	Alias  string `json:"alias"`
}

// Space is a Matrix space and its channels. Managed is true for a space linked
// to Registre, whose name and members Registre owns.
type Space struct {
	ID       int64     `json:"id"`
	RoomID   string    `json:"room_id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Managed  bool      `json:"managed"`
	Channels []Channel `json:"channels"`
}

// CreateSpaceRequest is the body of POST /spaces.
type CreateSpaceRequest struct {
	Name string `json:"name"`
}

// CreateChannelRequest is the body of POST /spaces/{id}/channels.
type CreateChannelRequest struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}
