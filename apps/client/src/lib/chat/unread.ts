export type RoomCounts = { total: number; highlight: number; muted: boolean; direct: boolean };
export type Badge = { count: number; strong: boolean; bold: boolean };

/**
 * What a room row shows. DMs count every message as a mention, channels only
 * highlights; a muted room keeps its highlights but loses the bold unread state.
 */
export function roomBadge(counts: RoomCounts): Badge {
	const count = counts.direct ? counts.total : counts.highlight;
	const bold = !counts.muted && counts.total > 0;
	return { count, strong: count > 0, bold: bold || count > 0 };
}

/** Sums rooms into one badge for a space or for the DM section. */
export function sumBadges(rooms: RoomCounts[]): Badge {
	return rooms.reduce<Badge>(
		(acc, room) => {
			const badge = roomBadge(room);
			return {
				count: acc.count + badge.count,
				strong: acc.strong || badge.strong,
				bold: acc.bold || badge.bold
			};
		},
		{ count: 0, strong: false, bold: false }
	);
}

export type Seen = { id: string; sender: string; message: boolean; mention: boolean };

/**
 * Counts unread messages from the newest event back to the user's read marker.
 * Synapse computes no notification counts for appservice users, and every
 * Prisme account is one, so the client counts from the loaded timeline.
 */
export function countUnread(events: Seen[], read: (id: string) => boolean, me: string) {
	let total = 0;
	let highlight = 0;
	for (let index = events.length - 1; index >= 0; index--) {
		const event = events[index];
		if (event.sender === me || read(event.id)) break;
		if (!event.message) continue;
		total++;
		if (event.mention) highlight++;
	}
	return { total, highlight };
}

/** Caps a counter for display. */
export function formatCount(count: number): string {
	return count > 99 ? '99+' : String(count);
}

/** Two messages share a header when one sender posts them within five minutes. */
export function continues(prev: { sender: string; ts: number } | undefined, next: { sender: string; ts: number }) {
	return !!prev && prev.sender === next.sender && next.ts - prev.ts < 5 * 60 * 1000;
}
