package users

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the directory route for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "users",
		Description: "Directory of people who have signed in to Prisme.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/users", Summary: "Search people by name, email or Matrix id",
				Description: "Query parameter q. Returns at most 20 entries.", Auth: "porte session", ResponseBody: []User{},
			},
		},
	}
}
