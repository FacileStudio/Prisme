import { KnownMembership, type MatrixClient, type Room } from 'matrix-js-sdk';
import type { Channel, Space } from '$lib/api/prisme';
import { isAgent } from '$lib/chat/users';
import { formatCount, roomBadge, sumBadges, type Badge, type RoomCounts } from '$lib/chat/unread';
import { dmPartner, dmRooms, isMuted, roomCounts, spaceInvites } from '$lib/matrix/rooms';
import { lastMessage, type Preview } from '$lib/matrix/preview';

export type ChannelRow = { roomId: string; name: string; badge: Badge; muted: boolean };
export type SpaceNav = { space: Space; name: string; joined: ChannelRow[]; joinable: ChannelRow[]; badge: Badge; counts: RoomCounts[] };
export type DmRow = {
	roomId: string;
	name: string;
	userId: string | null;
	invite: boolean;
	agent: boolean;
	badge: Badge;
	preview: Preview | null;
};
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

function dmRow(client: MatrixClient, room: Room, agents?: ReadonlyMap<string, unknown>): { row: DmRow; counts: RoomCounts } {
	const userId = dmPartner(client, room);
	const invite = room.getMyMembership() === KnownMembership.Invite;
	const member = userId ? room.getMember(userId) : null;
	const name = member?.rawDisplayName ?? client.getUser(userId ?? '')?.displayName ?? room.name;
	const counts = roomCounts(client, room, true);
	const badge = invite ? { count: 1, strong: true, bold: true } : roomBadge(counts);
	const agent = !!userId && isAgent(userId, agents);
	const preview = lastMessage(room, client.getUserId() ?? '');
	return { row: { roomId: room.roomId, name, userId, invite, agent, badge, preview }, counts };
}

/** Builds the sidebar model: spaces with their channels, DMs, and pending space invites. */
export function buildNav(client: MatrixClient, spaces: Space[], agents?: ReadonlyMap<string, unknown>): Nav {
	const dms = dmRooms(client).map((room) => dmRow(client, room, agents));
	return {
		spaces: spaces.map((space) => spaceNav(client, space)),
		dms: dms.map((d) => d.row),
		dmBadge: sumBadges(dms.map((d) => d.counts)),
		invites: spaceInvites(client)
	};
}

export type NavKind = 'messages' | 'browse' | 'channel';
export type NavPage = { label: string; href: string; active: boolean; group?: string; kind: NavKind };

export const BROWSE_PATH = '/chats';
export const MESSAGES_PATH = '/messages';

const roomPath = (roomId: string) => `/room/${encodeURIComponent(roomId)}`;

/** A row label carrying its unread state, the way the suite's rails show counters. */
export function countedLabel(name: string, badge: Badge): string {
	if (badge.count > 0) return `${name} · ${formatCount(badge.count)}`;
	return badge.bold ? `${name} ·` : name;
}

/**
 * The rail entries: Messages (every direct conversation, behind one entry),
 * the browse page, then the current space's joined channels.
 */
export function railPages(nav: Nav, current: SpaceNav | null, path: string, roomId: string | null): NavPage[] {
	const invites = nav.invites.length;
	const inDm = !!roomId && nav.dms.some((row) => row.roomId === roomId);
	const messages: NavPage = {
		label: countedLabel('Messages', nav.dmBadge),
		href: MESSAGES_PATH,
		active: path === MESSAGES_PATH || inDm,
		kind: 'messages'
	};
	const browse: NavPage = {
		label: invites > 0 ? `Overview · ${invites}` : 'Overview',
		href: BROWSE_PATH,
		active: path === BROWSE_PATH,
		kind: 'browse'
	};
	const channels: NavPage[] = (current?.joined ?? []).map((row) => ({
		label: countedLabel(row.name, row.badge),
		href: roomPath(row.roomId),
		active: row.roomId === roomId,
		group: current?.name ?? 'Channels',
		kind: 'channel'
	}));
	return [messages, browse, ...channels];
}

/** The space that lists a room as one of its channels. */
export function spaceOf(spaces: Space[], roomId: string): Space | undefined {
	return spaces.find((s) => s.channels.some((c) => c.room_id === roomId));
}
