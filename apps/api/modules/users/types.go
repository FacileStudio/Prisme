package users

// User is one directory entry: a person or an agent.
type User struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Agent     bool   `json:"agent"`
}

// Agent is one active agent, as GET /users/agents lists them.
type Agent struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	OwnerUserID string `json:"owner_user_id"`
}
