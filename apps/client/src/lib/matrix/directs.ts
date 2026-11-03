import { EventType, KnownMembership, Preset, type MatrixClient, type Room } from 'matrix-js-sdk';
import { directMap } from './rooms';

async function remember(client: MatrixClient, userId: string, roomId: string) {
	const map = directMap(client);
	const rooms = map[userId] ?? [];
	if (rooms.includes(roomId)) return;
	await client.setAccountData(EventType.Direct, { ...map, [userId]: [...rooms, roomId] } as never);
}

/** Opens the DM with a user, reusing a joined one before creating a new room. */
export async function openDm(client: MatrixClient, userId: string): Promise<string> {
	const existing = (directMap(client)[userId] ?? []).find(
		(id) => client.getRoom(id)?.getMyMembership() === KnownMembership.Join
	);
	if (existing) return existing;
	const { room_id } = await client.createRoom({
		is_direct: true,
		invite: [userId],
		preset: Preset.TrustedPrivateChat
	});
	await remember(client, userId, room_id);
	return room_id;
}

/** Accepts an invite; a DM invite is also recorded in `m.direct`. */
export async function acceptInvite(client: MatrixClient, room: Room, viaServer: string): Promise<void> {
	const inviter = room.getDMInviter();
	await client.joinRoom(room.roomId, { viaServers: [viaServer] });
	if (inviter && !room.isSpaceRoom()) await remember(client, inviter, room.roomId);
}
