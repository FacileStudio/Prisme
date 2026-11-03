package calls

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the call route for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "calls",
		Description: "Echo calls started from a Matrix room, and the call events Echo sends back.",
		Routes: []documentation.Route{
			{
				Method: "POST", Path: "/calls", Summary: "Start or join the call of a room",
				Description: "Gets or creates the Echo room, mints a LiveKit token and posts the call message unless join is true. " +
					"e2ee_key is the room's call encryption key, given only to members and never sent to Echo.",
				Auth: "porte session", RequestBody: Request{}, ResponseBody: Response{},
				Errors: callErrors(),
			},
			eventsRoute(),
		},
	}
}

func callErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 503, Code: "unavailable", Description: "Echo is not configured or did not answer."},
		{Status: 404, Code: "not_found", Description: "The caller is not in that room."},
	}
}

func eventsRoute() documentation.Route {
	return documentation.Route{
		Method: "POST", Path: EventsPath, Summary: "Receive an Echo call event",
		Description: "Echo posts participant_left and room_finished here, signed in X-Echo-Signature (hex HMAC-SHA256 of the body " +
			"with the integration key). When the last person leaves, the call message is edited into a finished call with its duration.",
		Auth:   "Echo signature",
		Errors: eventErrors(),
	}
}

func eventErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 401, Code: "unauthorized", Description: "The signature does not match."},
		{Status: 400, Code: "invalid", Description: "The body is not an Echo event."},
	}
}
