package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/FacileStudio/porte/oidc"
)

// UploadFromURL fetches an image over https and stores it in the media
// repository as the given user, returning its mxc URI. An mxc URI is returned
// as it is. The fetch goes through porte's guarded client, which refuses
// private addresses at connect time, because the URL comes from a caller.
func (c *Client) UploadFromURL(ctx context.Context, asUser, source string) (string, error) {
	if strings.HasPrefix(source, "mxc://") {
		return source, nil
	}
	data, contentType, err := oidc.FetchAvatar(ctx, source)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("matrix: avatar %s is %q, not an image", source, contentType)
	}
	return c.Upload(ctx, asUser, data, contentType)
}

// Upload stores bytes of the given content type in the media repository as
// the given user and returns the mxc URI.
func (c *Client) Upload(ctx context.Context, asUser string, data []byte, contentType string) (string, error) {
	query := url.Values{"user_id": {asUser}, "filename": {"avatar"}}
	target := c.cfg.BaseURL + "/_matrix/media/v3/upload?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.ASToken)
	req.Header.Set("Content-Type", contentType)
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", &Error{Status: resp.StatusCode, Code: "M_UPLOAD", Message: "avatar upload refused"}
	}
	var out struct {
		ContentURI string `json:"content_uri"`
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out.ContentURI, err
}
