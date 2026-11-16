package nuage

import (
	"context"
	"net/http"
	"strconv"
)

// Shared is a file shared by link: what a chat message needs to show it.
type Shared struct {
	FileID   int64  `json:"file_id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
}

type share struct {
	Token      string  `json:"token"`
	FileID     *int64  `json:"file_id"`
	Permission string  `json:"permission"`
	ExpiresAt  *string `json:"expires_at"`
}

type shareRequest struct {
	FileID     int64  `json:"file_id"`
	Permission string `json:"permission"`
}

// Share returns a view link to the file, reusing a lasting view link the user
// already made for it so a file shared twice keeps one link.
func (c *Client) Share(ctx context.Context, userID, fileID int64) (Shared, error) {
	var file File
	if err := c.get(ctx, userID, "/api/files/"+strconv.FormatInt(fileID, 10), &file); err != nil {
		return Shared{}, err
	}
	token, err := c.viewToken(ctx, userID, fileID)
	if err != nil {
		return Shared{}, err
	}
	return Shared{FileID: file.ID, Name: file.Name, MimeType: file.MimeType, Size: file.Size, URL: c.ShareURL(token)}, nil
}

func (c *Client) viewToken(ctx context.Context, userID, fileID int64) (string, error) {
	var mine struct {
		Shares []share `json:"shares"`
	}
	if err := c.get(ctx, userID, "/api/shares/by-me", &mine); err != nil {
		return "", err
	}
	for _, existing := range mine.Shares {
		if lasting(existing, fileID) {
			return existing.Token, nil
		}
	}
	var created share
	err := c.send(ctx, userID, call{method: http.MethodPost, path: "/api/shares", body: shareRequest{FileID: fileID, Permission: "view"}, out: &created})
	return created.Token, err
}

func lasting(s share, fileID int64) bool {
	return s.FileID != nil && *s.FileID == fileID && s.Permission == "view" && s.ExpiresAt == nil && s.Token != ""
}
