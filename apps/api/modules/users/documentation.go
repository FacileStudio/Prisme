package users

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the directory routes for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "users",
		Description: "Directory of the people and agents the caller shares a space with, and of its own agents.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/users", Summary: "Search the caller's space peers and agents by name, email or Matrix id",
				Description: "Query parameter q. Lists only the caller, the people and agents joined to a space the caller " +
					"is joined to, and the caller's own agents; nobody else, whatever q says. Returns at most 200 entries. " +
					"avatar_url is an mxc URI or empty; agent is true for an agent.",
				Auth: "porte session", ResponseBody: []User{},
			},
			{
				Method: "GET", Path: "/users/agents", Summary: "List the active agents the caller may see, with their owners",
				Description: "The caller's own agents, and the agents of space peers that are joined to a shared space. " +
					"The client draws agents as blobs and badges them from this list.",
				Auth: "porte session", ResponseBody: []Agent{},
			},
		},
	}
}
