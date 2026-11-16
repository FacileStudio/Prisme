package push

import "strings"

// TransactionsPath is where Synapse hands the appservice its events. The
// appservice spec fixes it; Traefik sends every public /_matrix path to
// Synapse, so only the compose network reaches this one.
const TransactionsPath = "/_matrix/app/v1/transactions/{txnId}"

// Config is the body of GET /push: Native is true when this server can push
// to the Android app.
type Config struct {
	Native bool `json:"native"`
}

// DeviceRequest is the body of POST /push/devices: the Firebase registration
// token of the Android app.
type DeviceRequest struct {
	Token string `json:"token"`
}

type transaction struct {
	Events []event `json:"events"`
}

type event struct {
	Type      string  `json:"type"`
	RoomID    string  `json:"room_id"`
	Sender    string  `json:"sender"`
	StateKey  *string `json:"state_key"`
	Timestamp int64   `json:"origin_server_ts"`
	Content   content `json:"content"`
}

type content struct {
	Body      string     `json:"body"`
	Mentions  *mentions  `json:"m.mentions"`
	RelatesTo *relation  `json:"m.relates_to"`
	Capsule   *struct{}  `json:"studio.facile.capsule"`
	Poll      *pollStart `json:"m.poll.start"`
	PollMSC   *pollStart `json:"org.matrix.msc3381.poll.start"`
}

type pollStart struct {
	Question pollText `json:"question"`
}

type pollText struct {
	Text    string `json:"m.text"`
	TextMSC string `json:"org.matrix.msc1767.text"`
	Body    string `json:"body"`
}

type mentions struct {
	UserIDs []string `json:"user_ids"`
}

type relation struct {
	RelType string `json:"rel_type"`
}

type device struct {
	Token        string
	MatrixUserID string
}

type message struct {
	Title  string
	Body   string
	RoomID string
}

// pollQuestion is the question of a poll start, stable or unstable, or empty
// when the content is not a poll.
func (c content) pollQuestion() string {
	for _, p := range []*pollStart{c.Poll, c.PollMSC} {
		if p == nil {
			continue
		}
		for _, text := range []string{p.Question.Text, p.Question.TextMSC, p.Question.Body} {
			if q := strings.TrimSpace(text); q != "" {
				return q
			}
		}
	}
	return ""
}
