// Package echo calls Echo's server-to-server integration API, which runs a
// LiveKit room per external id and mints join tokens for it.
package echo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

// State is who is connected to a room's latest call. Participants is empty
// when no call is running.
type State struct {
	RoomID       string        `json:"room_id"`
	Participants []Participant `json:"participants"`
}

// Event is one call event Echo forwards to the room's callback URL.
type Event struct {
	ID         string `json:"id"`
	Event      string `json:"event"`
	RoomID     string `json:"room_id"`
	ExternalID string `json:"external_id"`
	At         string `json:"at"`
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

// Room gets or creates the Echo room for an external id. A non-empty
// callback URL receives the room's call events.
func (c *Client) Room(ctx context.Context, externalID, title, callbackURL string) (Room, error) {
	var room Room
	body := map[string]string{"external_id": externalID, "title": title}
	if callbackURL != "" {
		body["callback_url"] = callbackURL
	}
	err := c.send(ctx, http.MethodPost, "/api/integrations/rooms", body, &room)
	return room, err
}

// State reads who is connected to an Echo room's call.
func (c *Client) State(ctx context.Context, roomID string) (State, error) {
	var state State
	err := c.send(ctx, http.MethodGet, "/api/integrations/rooms/"+url.PathEscape(roomID), nil, &state)
	return state, err
}

// Verify reports whether signature is Echo's hex HMAC-SHA256 of body, keyed
// with this app's integration secret.
func (c *Client) Verify(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(c.key))
	mac.Write(body)
	return hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(signature))
}

// Token mints a join token for a participant of an Echo room.
func (c *Client) Token(ctx context.Context, roomID string, who Participant) (Token, error) {
	var token Token
	err := c.send(ctx, http.MethodPost, "/api/integrations/rooms/"+url.PathEscape(roomID)+"/tokens", who, &token)
	return token, err
}

func (c *Client) send(ctx context.Context, method, path string, body, out any) error {
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
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
