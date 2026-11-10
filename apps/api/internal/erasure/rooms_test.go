package erasure

import (
	"encoding/json"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
)

func roomState(creator, version string, levels map[string]int64) []matrix.RoomEvent {
	create, _ := json.Marshal(map[string]string{"room_version": version})
	state := []matrix.RoomEvent{{Type: "m.room.create", Sender: creator, Content: create}}
	if levels != nil {
		content, _ := json.Marshal(map[string]any{"users": levels, "events": map[string]int{"m.room.power_levels": 100}})
		state = append(state, matrix.RoomEvent{Type: "m.room.power_levels", Sender: creator, Content: content})
	}
	return state
}

type roomCase struct {
	name    string
	members []string
	state   []matrix.RoomEvent
	tracked bool
	want    verdict
}

func roomCases() []roomCase {
	unsettled := "no member could administer the room before the erasure"
	return []roomCase{
		{"alone with their agent", []string{"@gone:hs", "@gone-bot:hs"}, roomState("@gone:hs", "12", nil), false,
			verdict{action: actionPurge}},
		{"a space with only the bot", []string{"@gone:hs", testBot}, roomState("@gone:hs", "12", nil), true,
			verdict{action: actionKeep}},
		{"a room with only the bot", []string{"@gone:hs", testBot}, roomState("@gone:hs", "12", nil), false,
			verdict{action: actionPurge}},
		{"v12, the other is the creator", []string{"@gone:hs", "@ada:hs"}, roomState("@ada:hs", "12", map[string]int64{}), false,
			verdict{action: actionKeep}},
		{"v12, the erased creator alone has power", []string{"@gone:hs", "@bo:hs", "@ada:hs"},
			roomState("@gone:hs", "12", map[string]int64{"@bo:hs": 50, "@ada:hs": 50}), false,
			verdict{action: actionPromote, promote: "@ada:hs"}},
		{"v10, the stronger member is promoted", []string{"@gone:hs", "@ada:hs", "@bo:hs"},
			roomState("@gone:hs", "10", map[string]int64{"@gone:hs": 100, "@bo:hs": 50}), false,
			verdict{action: actionPromote, promote: "@bo:hs"}},
		{"v10, another admin stays", []string{"@gone:hs", "@ada:hs"},
			roomState("@gone:hs", "10", map[string]int64{"@gone:hs": 100, "@ada:hs": 100}), false,
			verdict{action: actionKeep}},
		{"the bot administers it", []string{"@gone:hs", "@ada:hs", testBot},
			roomState("@gone:hs", "12", map[string]int64{testBot: 100}), true, verdict{action: actionKeep}},
		{"no admin before either", []string{"@gone:hs", "@ada:hs"},
			roomState("@zed:hs", "10", map[string]int64{"@zed:hs": 100}), false,
			verdict{action: actionKeep, unsettled: unsettled}},
		{"no power levels event, the creator is erased", []string{"@gone:hs", "@ada:hs"},
			roomState("@gone:hs", "10", nil), false, verdict{action: actionPromote, promote: "@ada:hs"}},
	}
}

func TestTheRoomRuleIsConservative(t *testing.T) {
	gone := map[string]bool{"@gone:hs": true, "@gone-bot:hs": true}
	for _, c := range roomCases() {
		got := decide(room{members: c.members, power: readPower(c.state), tracked: c.tracked}, gone, testBot)
		if got != c.want {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}
