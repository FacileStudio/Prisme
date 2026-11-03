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
				Auth: "porte session, or a delegated JWT with scope prisme.read", ResponseBody: []Space{},
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
		},
	}
}
