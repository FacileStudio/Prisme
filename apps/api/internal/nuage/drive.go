package nuage

import (
	"context"
	"net/url"
	"strconv"
)

// Space is a shared Nuage space the user belongs to.
type Space struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Folder is one folder in a listing.
type Folder struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}

// File is one file in a listing.
type File struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	MimeType  string `json:"mime_type"`
	Size      int64  `json:"size"`
	FolderID  *int64 `json:"folder_id"`
	UpdatedAt string `json:"updated_at"`
}

// Place is where a listing looks: a folder, a space's root, or the user's own
// root when both are nil.
type Place struct {
	FolderID *int64
	SpaceID  *int64
}

func (p Place) query(folderKey string) string {
	values := url.Values{}
	if p.FolderID != nil {
		values.Set(folderKey, strconv.FormatInt(*p.FolderID, 10))
	}
	if p.SpaceID != nil {
		values.Set("space_id", strconv.FormatInt(*p.SpaceID, 10))
	}
	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

// Spaces lists the shared spaces the user belongs to.
func (c *Client) Spaces(ctx context.Context, userID int64) ([]Space, error) {
	var answer struct {
		Spaces []Space `json:"spaces"`
	}
	err := c.get(ctx, userID, "/api/spaces", &answer)
	return answer.Spaces, err
}

// List returns the folders and files directly inside place.
func (c *Client) List(ctx context.Context, userID int64, place Place) ([]Folder, []File, error) {
	var folders struct {
		Folders []Folder `json:"folders"`
	}
	if err := c.get(ctx, userID, "/api/folders"+place.query("parent_id"), &folders); err != nil {
		return nil, nil, err
	}
	var files struct {
		Files []File `json:"files"`
	}
	err := c.get(ctx, userID, "/api/files"+place.query("folder_id"), &files)
	return folders.Folders, files.Files, err
}

// Search finds the user's files whose name matches query.
func (c *Client) Search(ctx context.Context, userID int64, query string) ([]File, error) {
	var answer struct {
		Results []struct {
			File
			Type string `json:"type"`
		} `json:"results"`
	}
	path := "/api/search?" + url.Values{"q": {query}, "type": {"file"}, "limit": {"40"}}.Encode()
	if err := c.get(ctx, userID, path, &answer); err != nil {
		return nil, err
	}
	files := make([]File, 0, len(answer.Results))
	for _, result := range answer.Results {
		if result.Type == "file" {
			files = append(files, result.File)
		}
	}
	return files, nil
}
