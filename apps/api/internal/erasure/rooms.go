package erasure

import "slices"

type action int

const (
	actionKeep action = iota
	actionPurge
	actionPromote
)

// room is what the rule reads about one room: who is joined, who holds
// power, and whether Prisme lists it as a space or a channel.
type room struct {
	members []string
	power   power
	tracked bool
}

// verdict is the rule's answer: the action, the member to promote, and, for
// a room kept as it was although it may need an operator, why.
type verdict struct {
	action    action
	promote   string
	unsettled string
}

// decide is the conservative room rule. A room nobody else is in is purged,
// unless it is a Prisme space or channel, which stays with the bot. A shared
// room stays for its other members; if none of them, nor the bot, could
// administer it after the erasure while an erased account could before, the
// most powerful remaining person is raised to that account's power.
func decide(r room, erased map[string]bool, bot string) verdict {
	others := slices.DeleteFunc(slices.Clone(r.members), func(m string) bool { return erased[m] || m == bot })
	if len(others) == 0 {
		if r.tracked {
			return verdict{action: actionKeep}
		}
		return verdict{action: actionPurge}
	}
	if slices.ContainsFunc(r.members, func(m string) bool { return !erased[m] && r.power.admin(m) }) {
		return verdict{action: actionKeep}
	}
	if !slices.ContainsFunc(r.members, func(m string) bool { return erased[m] && r.power.admin(m) }) {
		return verdict{action: actionKeep, unsettled: "no member could administer the room before the erasure"}
	}
	return verdict{action: actionPromote, promote: strongest(others, r.power)}
}

// strongest is the member with the highest power, the lowest id on a tie.
func strongest(members []string, p power) string {
	sorted := slices.Clone(members)
	slices.Sort(sorted)
	best := sorted[0]
	for _, m := range sorted[1:] {
		if p.level(m) > p.level(best) {
			best = m
		}
	}
	return best
}
