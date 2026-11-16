package nuage

import "github.com/FacileStudio/Prisme/apps/api/internal/nuage"

// SpacesResponse lists the user's shared Nuage spaces.
type SpacesResponse struct {
	Spaces []nuage.Space `json:"spaces"`
}

// ListResponse is the content of one folder, space or the user's root.
type ListResponse struct {
	Folders []nuage.Folder `json:"folders"`
	Files   []nuage.File   `json:"files"`
}

// SearchResponse is the user's files matching a name.
type SearchResponse struct {
	Files []nuage.File `json:"files"`
}

// ShareRequest names the file to share into a chat.
type ShareRequest struct {
	FileID int64 `json:"file_id"`
}

// ShareResponse is the view link and what a chat message shows of the file.
type ShareResponse = nuage.Shared
