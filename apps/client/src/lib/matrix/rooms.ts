import { EventType, KnownMembership, NotificationCountType, type MatrixClient, type Room } from 'matrix-js-sdk';
import type { RoomCounts } from '$lib/chat/unread';

type DirectMap = Record<string, string[]>;

/** The `m.direct` account data: partner user id to DM room ids. */
export function directMap(client: MatrixClient): DirectMap {
	return (client.getAccountData(EventType.Direct)?.getContent() ?? {}) as DirectMap;
}

function directRoomIds(client: MatrixClient): Set<string> {
	return new Set(Object.values(directMap(client)).flat());
}

/** Reports whether a room is a DM, from `m.direct` or the invite's `is_direct` flag. */
export function isDirect(client: MatrixClient, room: Room): boolean {
	if (directRoomIds(client).has(room.roomId)) return true;
	return room.getMyMembership() === KnownMembership.Invite && !!room.getDMInviter() && !room.isSpaceRoom();
}

/** Joined and invited DM rooms, most recently active first. */
export function dmRooms(client: MatrixClient): Room[] {
	const ids = directRoomIds(client);
	return client
		.getRooms()
		.filter((room) => {
			const membership = room.getMyMembership();
			if (membership !== KnownMembership.Join && membership !== KnownMembership.Invite) return false;
			return ids.has(room.roomId) || isDirect(client, room);
		})
		.sort((a, b) => b.getLastActiveTimestamp() - a.getLastActiveTimestamp());
}

/** The other person in a DM. */
export function dmPartner(client: MatrixClient, room: Room): string | null {
	const entry = Object.entries(directMap(client)).find(([, rooms]) => rooms.includes(room.roomId));
	if (entry) return entry[0];
	const me = client.getUserId();
	const other = room.getMembers().find((m) => m.userId !== me);
	return other?.userId ?? room.getDMInviter() ?? null;
}

/** Reports whether the room has a muting room push rule. */
export function isMuted(client: MatrixClient, roomId: string): boolean {
	const rule = client.getRoomPushRule('global', roomId);
	if (!rule) return false;
	return rule.actions.length === 0 || rule.actions.includes('dont_notify' as never);
}

/** The unread counters of one room. */
export function roomCounts(client: MatrixClient, room: Room, direct: boolean): RoomCounts {
	return {
		total: room.getUnreadNotificationCount(NotificationCountType.Total),
		highlight: room.getUnreadNotificationCount(NotificationCountType.Highlight),
		muted: isMuted(client, room.roomId),
		direct
	};
}

/** Spaces the user is invited to but has not joined yet. */
export function spaceInvites(client: MatrixClient): Room[] {
	return client.getRooms().filter((r) => r.getMyMembership() === KnownMembership.Invite && r.isSpaceRoom());
}
