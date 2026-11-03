package agents

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

const agentAuth = "porte session, or a Registre token-exchange JWT (aud prisme) with scope prisme.agents"

// Documentation describes the agent routes for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "agents",
		Description: "Bot accounts (@agent-<owner>-<name>) a user owns, for Bulle assistants.",
		Routes: []documentation.Route{
			{Method: "GET", Path: "/agents", Summary: "List the caller's agents and their rooms", Auth: agentAuth, ResponseBody: []Agent{}},
			{
				Method: "POST", Path: "/agents", Summary: "Create or return the caller's agent with this name",
				Description: "Idempotent per (owner, name). Sets the display name and avatar every call.",
				Auth:        agentAuth, RequestBody: CreateRequest{}, ResponseBody: Login{},
			},
			{
				Method: "PUT", Path: "/agents/{user_id}/rooms", Summary: "Set the rooms the agent is joined to",
				Description: "Invites and joins the listed rooms, leaves the others. Only rooms the owner is joined to.",
				Auth:        agentAuth, RequestBody: RoomsRequest{}, ResponseBody: RoomsResponse{},
				Errors: []documentation.Error{{Status: 403, Code: "permission_denied", Description: "A room the owner is not in."}},
			},
			{Method: "DELETE", Path: "/agents/{user_id}", Summary: "Deactivate the agent", Auth: agentAuth},
		},
	}
}
