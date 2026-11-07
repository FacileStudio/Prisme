package push

import (
	"bytes"
	"cmp"
	"context"
	"encoding/base64"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2/jwt"
)

const (
	fcmScope    = "https://www.googleapis.com/auth/firebase.messaging"
	tokenURL    = "https://oauth2.googleapis.com/token"
	sendTimeout = 10 * time.Second
)

// errGone means Firebase no longer knows the device, so its pusher can go.
var errGone = stderrors.New("push device is gone")

// FCM sends Android notifications through Firebase Cloud Messaging HTTP v1.
type FCM struct {
	endpoint string
	client   *http.Client
}

type serviceAccount struct {
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

// NewFCM builds a sender from a Firebase service account, given as raw JSON or
// base64 JSON. An empty value returns nil: push stays off.
func NewFCM(raw string) (*FCM, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if !strings.HasPrefix(raw, "{") {
		decoded, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return nil, fmt.Errorf("FCM_SERVICE_ACCOUNT is neither JSON nor base64: %w", err)
		}
		raw = string(decoded)
	}
	var account serviceAccount
	if err := json.Unmarshal([]byte(raw), &account); err != nil {
		return nil, fmt.Errorf("read FCM_SERVICE_ACCOUNT: %w", err)
	}
	if account.ProjectID == "" || account.ClientEmail == "" || account.PrivateKey == "" {
		return nil, stderrors.New("FCM_SERVICE_ACCOUNT lacks project_id, client_email or private_key")
	}
	cfg := &jwt.Config{
		Email: account.ClientEmail, PrivateKey: []byte(account.PrivateKey),
		Scopes: []string{fcmScope}, TokenURL: cmp.Or(account.TokenURI, tokenURL),
	}
	client := cfg.Client(context.Background())
	client.Timeout = sendTimeout
	return &FCM{endpoint: "https://fcm.googleapis.com/v1/projects/" + account.ProjectID + "/messages:send", client: client}, nil
}

type fcmRequest struct {
	Message fcmMessage `json:"message"`
}

type fcmMessage struct {
	Token        string            `json:"token"`
	Notification map[string]string `json:"notification"`
	Data         map[string]string `json:"data"`
	Android      fcmAndroid        `json:"android"`
}

type fcmAndroid struct {
	Priority     string            `json:"priority"`
	Notification map[string]string `json:"notification"`
}

func (f *FCM) send(ctx context.Context, token string, m message) error {
	body, err := json.Marshal(fcmRequest{Message: fcmMessage{
		Token:        token,
		Notification: map[string]string{"title": m.Title, "body": m.Body},
		Data:         map[string]string{"room_id": m.RoomID},
		Android:      fcmAndroid{Priority: "HIGH", Notification: map[string]string{"tag": m.RoomID, "icon": "ic_stat_prisme"}},
	}})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	reply, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode == http.StatusNotFound || bytes.Contains(reply, []byte("UNREGISTERED")) {
		return errGone
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("fcm answered %d", resp.StatusCode)
	}
	return nil
}
