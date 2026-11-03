package session

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the session route for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "session",
		Description: "The signed-in user and their Matrix login.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/session", Summary: "Return the user and a Matrix access token",
				Description: "Creates the Matrix account on first call. One device per porte session, revoked on logout.",
				Auth:        "porte session", ResponseBody: Response{},
			},
		},
	}
}
