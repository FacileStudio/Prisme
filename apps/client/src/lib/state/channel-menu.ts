import { goto } from '$app/navigation';
import { copyText, icons, toast, type MenuEntry } from '@facile/muse';
import { glyphs } from '$lib/glyphs';
import { isMuted } from '$lib/matrix/rooms';
import { canRename } from '$lib/matrix/power';
import { t } from '$lib/i18n/index.svelte';
import { app } from './app.svelte';
import { joinChannel, roomHref, setMuted } from './actions';
import { ui } from './ui.svelte';

async function attempt(work: () => Promise<unknown>, done: string, failed: string) {
	try {
		await work();
		if (done) toast.success(done);
	} catch (err) {
		toast.danger(err instanceof Error ? err.message : failed);
	}
}

async function markRead(roomId: string) {
	const room = app.client?.getRoom(roomId);
	const last = room ? [...room.getLiveTimeline().getEvents()].reverse().find((e) => e.getId()?.startsWith('$')) : undefined;
	if (last) await app.client?.sendReadReceipt(last);
}

async function copyLink(roomId: string) {
	const copied = await copyText(new URL(roomHref(roomId), location.origin).href);
	if (copied) toast.success(t('channel.copied'));
	else toast.danger(t('channel.copyFailed'));
}

function joinedEntries(roomId: string): MenuEntry[] {
	const client = app.client;
	const room = client?.getRoom(roomId);
	if (!client || !room) return [];
	const muted = isMuted(client, roomId);
	const entries: MenuEntry[] = [
		{ label: t('channel.open'), icon: icons.arrow, onselect: () => void goto(roomHref(roomId)) },
		{ label: t('channel.markRead'), icon: icons.check, onselect: () => void attempt(() => markRead(roomId), t('channel.markedRead'), t('channel.markRead')) },
		{
			label: t(muted ? 'menu.unmuteChannel' : 'menu.muteChannel'),
			icon: muted ? glyphs.bell : glyphs.bellOff,
			onselect: () => void attempt(() => setMuted(roomId, !muted), t(muted ? 'menu.unmuted' : 'menu.muted'), t('menu.muteFailed'))
		},
		{ label: t('channel.copyLink'), icon: icons.copy, onselect: () => void copyLink(roomId) }
	];
	if (canRename(client, room)) {
		entries.push({ separator: true }, { label: t('menu.renameChannel'), icon: icons.edit, onselect: () => ui.showFor('rename-room', roomId) });
	}
	entries.push({ separator: true }, { label: t('channel.leave'), icon: glyphs.exit, tone: 'danger', onselect: () => ui.showFor('leave-room', roomId) });
	return entries;
}

function joinableEntries(roomId: string): MenuEntry[] {
	return [
		{
			label: t('channel.join'),
			icon: glyphs.login,
			onselect: () => void attempt(async () => (await joinChannel(roomId), goto(roomHref(roomId))), '', t('space.joinFailed'))
		},
		{ label: t('channel.copyLink'), icon: icons.copy, onselect: () => void copyLink(roomId) }
	];
}

/** The right-click entries for a channel: joined channels get the full set, joinable ones only Join and Copy link. */
export function channelMenu(roomId: string): MenuEntry[] {
	const membership = app.client?.getRoom(roomId)?.getMyMembership();
	return membership === 'join' ? joinedEntries(roomId) : joinableEntries(roomId);
}

/** The channel menu for the row under a right-clicked element, read from its `data-room-id`. */
export function channelMenuAt(target: Element): MenuEntry[] {
	const roomId = target.closest<HTMLElement>('[data-room-id]')?.dataset.roomId;
	return roomId ? channelMenu(roomId) : [];
}
