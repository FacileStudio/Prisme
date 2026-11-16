package nuage

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the Nuage picker for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name: "nuage",
		Description: "A picker over the user's Nuage drive. Prisme exchanges the user's Registre token for Nuage and acts " +
			"as them, so Nuage's own access rules decide what shows. A shared file stays hosted on Nuage.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/nuage/spaces", Summary: "List Nuage spaces",
				Description: "The shared spaces the user belongs to in Nuage.",
				Auth:        "porte session", ResponseBody: SpacesResponse{}, Errors: pickerErrors(),
			},
			{
				Method: "GET", Path: "/nuage/list", Summary: "List a folder",
				Description: "Query parameters folder_id and space_id, both optional. Neither lists the user's own root.",
				Auth:        "porte session", ResponseBody: ListResponse{}, Errors: pickerErrors(),
			},
			{
				Method: "GET", Path: "/nuage/search", Summary: "Search files",
				Description: "Query parameter q, 1 to 200 characters. Matches file names across the user's drive.",
				Auth:        "porte session", ResponseBody: SearchResponse{}, Errors: pickerErrors(),
			},
			{
				Method: "POST", Path: "/nuage/share", Summary: "Share a file into a chat",
				Description: "Returns a view link to the file, reusing a lasting view link the user already made for it.",
				Auth:        "porte session", RequestBody: ShareRequest{}, ResponseBody: ShareResponse{}, Errors: pickerErrors(),
			},
		},
	}
}

func pickerErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 400, Code: "invalid_argument", Description: "A malformed id or query."},
		{Status: 403, Code: "permission_denied", Description: "Nuage refused Prisme's delegated access."},
		{Status: 404, Code: "not_found", Description: "The file or folder is not in the user's Nuage."},
		{Status: 412, Code: "failed_precondition", Description: "The user must sign in again before Prisme can reach Nuage."},
		{Status: 503, Code: "unavailable", Description: "Nuage is not connected or did not answer."},
	}
}
