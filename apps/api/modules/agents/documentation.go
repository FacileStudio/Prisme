package agents

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

const agentAuth = "porte session, or a Registre token-exchange JWT (aud prisme) with scope prisme.agents"

// Documentation describes the agent routes for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "agents",
		Description: "Bot accounts (@name) a user owns, for Bulle assistants.",
		Routes: []documentation.Route{
			{Method: "GET", Path: "/agents", Summary: "List the caller's agents and their rooms", Auth: agentAuth, ResponseBody: []Agent{}},
			{
				Method: "POST", Path: "/agents", Summary: "Create or return the caller's agent with this name",
				Description: "Idempotent per (owner, name), case-insensitive. The agent is @<slug>, or @<slug>-<owner localpart> " +
					"when another owner or a person holds the slug. A legacy @agent-<owner>-<slug> agent moves to that id: " +
					"the new account joins the old one's rooms and the old account is deactivated. Sets the display name " +
					"every call; an empty avatar_url gives the agent the blob its name draws.",
				Auth: agentAuth, RequestBody: CreateRequest{}, ResponseBody: Login{},
			},
			{
				Method: "GET", Path: "/agents/check", Summary: "Check whether an agent name is available",
				Description: "slug is the Matrix localpart POST /agents would give the agent.",
				Auth:        agentAuth, ResponseBody: AvailabilityResponse{},
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
