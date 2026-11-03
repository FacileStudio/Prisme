package calls

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the call route for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "calls",
		Description: "Echo calls started from a Matrix room.",
		Routes: []documentation.Route{
			{
				Method: "POST", Path: "/calls", Summary: "Start or join the call of a room",
				Description: "Gets or creates the Echo room, mints a LiveKit token and posts the call message unless join is true.",
				Auth:        "porte session", RequestBody: Request{}, ResponseBody: Response{},
				Errors: callErrors(),
			},
		},
	}
}

func callErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 503, Code: "unavailable", Description: "Echo is not configured or did not answer."},
		{Status: 404, Code: "not_found", Description: "The caller is not in that room."},
	}
}
