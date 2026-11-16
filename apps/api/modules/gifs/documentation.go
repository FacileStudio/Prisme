package gifs

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

var unavailable = documentation.Error{Status: 503, Code: "unavailable", Description: "No GIPHY key, or GIPHY did not answer."}

var searchRoute = documentation.Route{
	Method: "GET", Path: "/gifs", Summary: "Search GIFs",
	Description: "q searches, an empty q lists trending GIFs. lang is en, fr, es or de; offset pages. Rated G.",
	Auth:        "porte session", ResponseBody: PageResponse{},
	Errors: []documentation.Error{
		{Status: 400, Code: "invalid_argument", Description: "q is longer than 100 characters."},
		unavailable,
	},
}

var mediaRoute = documentation.Route{
	Method: "GET", Path: "/gifs/media", Summary: "Fetch a GIF",
	Description: "u is a GIPHY media URL from a search result. Only https links on GIPHY's media hosts are fetched.",
	Auth:        "porte session",
	Errors: []documentation.Error{
		{Status: 400, Code: "invalid_argument", Description: "u is not a GIPHY media URL."},
		unavailable,
	},
}

// Documentation describes the GIF picker's routes for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name: "gifs",
		Description: "The GIF picker, powered by GIPHY. Prisme searches and fetches from the server, so the browser never " +
			"talks to GIPHY. Off, with 503, while GIPHY_API_KEY is unset.",
		Routes: []documentation.Route{searchRoute, mediaRoute},
	}
}
