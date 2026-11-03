// Package echo calls Echo's server-to-server integration API, which runs a
// LiveKit room per external id and mints join tokens for it.
package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an Echo integration client authenticated by the app's secret.
type Client struct {
	baseURL string
	key     string
	http    *http.Client
}

// Room is an Echo room bound to an external id.
type Room struct {
	RoomID  string `json:"room_id"`
	JoinURL string `json:"join_url"`
}

// Token is a LiveKit join token for one participant.
type Token struct {
	LiveKitURL string `json:"livekit_url"`
	Token      string `json:"token"`
}

// Participant is who asks for a token.
type Participant struct {
	Identity  string `json:"identity"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// New builds a client for the Echo deployment at baseURL.
func New(baseURL, key string) *Client {
	return &Client{baseURL: baseURL, key: key, http: &http.Client{Timeout: 15 * time.Second}}
}

// Room gets or creates the Echo room for an external id.
func (c *Client) Room(ctx context.Context, externalID, title string) (Room, error) {
	var room Room
	body := map[string]string{"external_id": externalID, "title": title}
	err := c.post(ctx, "/api/integrations/rooms", body, &room)
	return room, err
}

// Token mints a join token for a participant of an Echo room.
func (c *Client) Token(ctx context.Context, roomID string, who Participant) (Token, error) {
	var token Token
	err := c.post(ctx, "/api/integrations/rooms/"+url.PathEscape(roomID)+"/tokens", who, &token)
	return token, err
}

func (c *Client) post(ctx context.Context, path string, body, out any) error {
	encoded, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("echo: %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("echo: %s answered %d: %s", path, resp.StatusCode, bytes.TrimSpace(data))
	}
	return json.Unmarshal(data, out)
}
