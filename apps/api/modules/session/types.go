package session

// User is the signed-in person as the client shows them.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AvatarURL    string `json:"avatar_url"`
	MatrixUserID string `json:"matrix_user_id"`
}

// MatrixLogin is what a Matrix client needs to start: the public homeserver
// URL and a device token.
type MatrixLogin struct {
	Homeserver  string `json:"homeserver"`
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
}

// Response is the body of GET /session.
type Response struct {
	User         User        `json:"user"`
	Matrix       MatrixLogin `json:"matrix"`
	CallsEnabled bool        `json:"calls_enabled"`
}
