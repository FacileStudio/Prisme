import type { MatrixClient, Room } from 'matrix-js-sdk';
import { prisme, type Space } from '$lib/api/prisme';
import { acceptInvite } from '$lib/matrix/directs';
import { t } from '$lib/i18n/index.svelte';
import { app } from './app.svelte';

/** The link to a room's page. */
export function roomHref(roomId: string): string {
	return `/room/${encodeURIComponent(roomId)}`;
}

/** The homeserver name, taken from the signed-in user's id. */
export function serverName(): string {
	return app.session?.matrix.user_id.split(':').slice(1).join(':') ?? '';
}

function client(): MatrixClient {
	if (!app.client) throw new Error(t('app.clientNotRunning'));
	return app.client;
}

/** Joins a channel of a space the user is a member of (restricted join rule). */
export async function joinChannel(roomId: string): Promise<void> {
	await client().joinRoom(roomId, { viaServers: [serverName()] });
}

/** Accepts a space invite, then joins the space's #general. */
export async function acceptSpace(room: Room): Promise<Space | undefined> {
	await acceptInvite(client(), room, serverName());
	await app.refreshSpaces();
	const space = app.spaces.find((s) => s.room_id === room.roomId);
	const general = space?.channels.find((c) => c.name === 'general') ?? space?.channels[0];
	if (general) await joinChannel(general.room_id);
	return space;
}

/** Creates a space and reloads the list. */
export async function createSpace(name: string): Promise<Space> {
	const space = await prisme.createSpace(name);
	await app.refreshSpaces();
	return space;
}

/** Creates a channel in a space and reloads the list. */
export async function createChannel(spaceId: number, name: string, topic: string): Promise<string> {
	const channel = await prisme.createChannel(spaceId, name, topic);
	await app.refreshSpaces();
	return channel.room_id;
}

/** Mutes or unmutes a room with a room push rule. */
export async function setMuted(roomId: string, muted: boolean): Promise<void> {
	await client().setRoomMutePushRule('global', roomId, muted);
}
