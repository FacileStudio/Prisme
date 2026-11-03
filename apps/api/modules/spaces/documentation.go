package spaces

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the space routes for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name:        "spaces",
		Description: "Matrix spaces and their channels, provisioned through the appservice.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/spaces", Summary: "List the spaces the caller is joined to",
				Description: "managed is true for a space linked to Registre, whose name and members Registre owns.",
				Auth:        "porte session, or a delegated JWT with scope prisme.read", ResponseBody: []Space{},
			},
			{
				Method: "POST", Path: "/spaces", Summary: "Create a space with a #general channel",
				Auth: "porte session", RequestBody: CreateSpaceRequest{}, ResponseBody: Space{},
			},
			{
				Method: "POST", Path: "/spaces/{id}/channels", Summary: "Create a channel in a space",
				Auth: "porte session", RequestBody: CreateChannelRequest{}, ResponseBody: Channel{},
				Errors: []documentation.Error{{Status: 404, Code: "not_found", Description: "No such space, or the caller is not in it."}},
			},
			{
				Method: "POST", Path: "/registre/sync", Summary: "Nudge a pull of Registre's shared spaces",
				Description: "Registre calls it after every change; the body is ignored and the answer is 202. " +
					"Mounted only when OIDC_ISSUER, OIDC_CLIENT_ID and OIDC_CLIENT_SECRET are set.",
				Auth: "none",
			},
		},
	}
}
