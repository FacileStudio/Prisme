package unfurl

import documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"

// Documentation describes the link preview route for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name: "unfurl",
		Description: "Previews of links to other Facile apps. Prisme fetches only https links to Feuilles, Docs, Echo, " +
			"Clap, Diapo and Nuage, without cookies, and never Capsule links, which a fetch would burn.",
		Routes: []documentation.Route{
			{
				Method: "GET", Path: "/unfurl", Summary: "Preview a suite link",
				Description: "Query parameter url. Reads the app's oEmbed answer, else its Open Graph tags. A link the app " +
					"keeps private answers a card with app alone. The fragment is dropped before any fetch.",
				Auth: "porte session", ResponseBody: Preview{},
				Errors: previewErrors(),
			},
		},
	}
}

func previewErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 400, Code: "invalid_argument", Description: "Not an https link to an app Prisme previews."},
	}
}
