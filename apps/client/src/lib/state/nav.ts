import { KnownMembership, type MatrixClient, type Room } from 'matrix-js-sdk';
import type { Channel, Space } from '$lib/api/prisme';
import { isAgent } from '$lib/chat/users';
import { roomBadge, sumBadges, type Badge, type RoomCounts } from '$lib/chat/unread';
import { dmPartner, dmRooms, isMuted, roomCounts, spaceInvites } from '$lib/matrix/rooms';

export type ChannelRow = { roomId: string; name: string; badge: Badge; muted: boolean };
export type SpaceNav = { space: Space; name: string; joined: ChannelRow[]; joinable: ChannelRow[]; badge: Badge; counts: RoomCounts[] };
export type DmRow = { roomId: string; name: string; userId: string | null; invite: boolean; agent: boolean; badge: Badge };
export type Nav = { spaces: SpaceNav[]; dms: DmRow[]; dmBadge: Badge; invites: Room[] };

const NONE: Badge = { count: 0, strong: false, bold: false };

function channelRow(client: MatrixClient, channel: Channel): { row: ChannelRow; counts: RoomCounts | null } {
	const room = client.getRoom(channel.room_id);
	const joined = room?.getMyMembership() === KnownMembership.Join;
	const counts = joined && room ? roomCounts(client, room, false) : null;
	const name = joined && room ? room.name : channel.name;
	const row = { roomId: channel.room_id, name, badge: counts ? roomBadge(counts) : NONE, muted: isMuted(client, channel.room_id) };
	return { row, counts };
}

function spaceNav(client: MatrixClient, space: Space): SpaceNav {
	const rows = [...space.channels].sort((a, b) => a.name.localeCompare(b.name)).map((c) => channelRow(client, c));
	const counts = rows.flatMap((r) => (r.counts ? [r.counts] : []));
	return {
		space,
		name: client.getRoom(space.room_id)?.name ?? space.name,
		joined: rows.filter((r) => r.counts).map((r) => r.row),
		joinable: rows.filter((r) => !r.counts).map((r) => r.row),
		badge: sumBadges(counts),
		counts
	};
}

function dmRow(client: MatrixClient, room: Room): { row: DmRow; counts: RoomCounts } {
	const userId = dmPartner(client, room);
	const invite = room.getMyMembership() === KnownMembership.Invite;
	const member = userId ? room.getMember(userId) : null;
	const name = member?.name ?? client.getUser(userId ?? '')?.displayName ?? room.name;
	const counts = roomCounts(client, room, true);
	const badge = invite ? { count: 1, strong: true, bold: true } : roomBadge(counts);
	return { row: { roomId: room.roomId, name, userId, invite, agent: !!userId && isAgent(userId), badge }, counts };
}

/** Builds the sidebar model: spaces with their channels, DMs, and pending space invites. */
export function buildNav(client: MatrixClient, spaces: Space[]): Nav {
	const dms = dmRooms(client).map((room) => dmRow(client, room));
	return {
		spaces: spaces.map((space) => spaceNav(client, space)),
		dms: dms.map((d) => d.row),
		dmBadge: sumBadges(dms.map((d) => d.counts)),
		invites: spaceInvites(client)
	};
}

/** The space that lists a room as one of its channels. */
export function spaceOf(spaces: Space[], roomId: string): Space | undefined {
	return spaces.find((s) => s.channels.some((c) => c.room_id === roomId));
}
