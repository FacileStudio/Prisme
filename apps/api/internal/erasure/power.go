package erasure

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
)

const (
	defaultStateLevel = 50
	creatorLevel      = 100
	creatorPowerFrom  = 12
)

// power is who can change a room's power levels, which is what administering
// a room comes down to. From room version 12 a creator holds unlimited power
// and is listed in no power levels event.
type power struct {
	creators     map[string]bool
	creatorPower bool
	levels       bool
	users        map[string]int64
	usersDefault int64
	required     int64
}

type levelsContent struct {
	Users        map[string]int64 `json:"users"`
	UsersDefault int64            `json:"users_default"`
	StateDefault *int64           `json:"state_default"`
	Events       map[string]int64 `json:"events"`
}

type createContent struct {
	RoomVersion        string   `json:"room_version"`
	AdditionalCreators []string `json:"additional_creators"`
}

func readPower(state []matrix.RoomEvent) power {
	p := power{creators: map[string]bool{}, users: map[string]int64{}, required: defaultStateLevel}
	for _, event := range state {
		if event.StateKey != "" {
			continue
		}
		switch event.Type {
		case "m.room.create":
			p.readCreate(event)
		case "m.room.power_levels":
			p.readLevels(event.Content)
		}
	}
	return p
}

func (p *power) readCreate(event matrix.RoomEvent) {
	p.creators[event.Sender] = true
	var content createContent
	if json.Unmarshal(event.Content, &content) != nil {
		return
	}
	version, err := strconv.Atoi(content.RoomVersion)
	p.creatorPower = err == nil && version >= creatorPowerFrom
	if p.creatorPower {
		for _, creator := range content.AdditionalCreators {
			p.creators[creator] = true
		}
	}
}

func (p *power) readLevels(raw json.RawMessage) {
	var content levelsContent
	if json.Unmarshal(raw, &content) != nil {
		return
	}
	p.levels = true
	p.users, p.usersDefault = content.Users, content.UsersDefault
	if content.StateDefault != nil {
		p.required = *content.StateDefault
	}
	if level, ok := content.Events["m.room.power_levels"]; ok {
		p.required = level
	}
}

func (p power) level(user string) int64 {
	switch {
	case p.creators[user] && p.creatorPower:
		return math.MaxInt64
	case !p.levels && p.creators[user]:
		return creatorLevel
	case !p.levels:
		return 0
	}
	if level, ok := p.users[user]; ok {
		return level
	}
	return p.usersDefault
}

func (p power) admin(user string) bool {
	if !p.levels {
		return p.creators[user]
	}
	return p.level(user) >= p.required
}
