package push

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the push routes for the API reference. The
// appservice routes Synapse calls live outside /api, under /_matrix/app/v1, as
// the Matrix spec fixes.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "push",
		Description: "Notifications for the Android app: direct messages and mentions, sent through Firebase.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/push", Summary: "Say whether this server pushes to the Android app",
				Auth: "porte session", ResponseBody: Config{},
			},
			{
				Method: "POST", Path: "/push/devices", Summary: "Register the Android app's Firebase token",
				Description: "Stores the token for this session, so signing out removes it. " +
					"Direct messages and mentions in rooms the user has not muted then reach the phone.",
				Auth: "porte session", RequestBody: DeviceRequest{}, Errors: deviceErrors(),
			},
		},
	}
}

func deviceErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 503, Code: "unavailable", Description: "FCM_SERVICE_ACCOUNT or PRISME_HS_TOKEN is not set."},
		{Status: 400, Code: "invalid_argument", Description: "The token is empty or too long."},
	}
}
