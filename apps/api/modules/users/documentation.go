package users

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the directory routes for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "users",
		Description: "Directory of people who have signed in to Prisme, and of their agents.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/users", Summary: "Search people and agents by name, email or Matrix id",
				Description: "Query parameter q. Returns at most 20 entries. avatar_url is an mxc URI or empty; " +
					"agent is true for an agent.",
				Auth: "porte session", ResponseBody: []User{},
			},
			{
				Method: "GET", Path: "/users/agents", Summary: "List every active agent and its owner",
				Description: "The client draws agents as blobs and badges them from this list.",
				Auth:        "porte session", ResponseBody: []Agent{},
			},
		},
	}
}
