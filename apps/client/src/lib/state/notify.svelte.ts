import { RoomEvent, type MatrixClient, type MatrixEvent, type Room } from 'matrix-js-sdk';
import { isDirect, isMuted } from '$lib/matrix/rooms';
import { isMessage } from '$lib/matrix/timeline';

const STORAGE_KEY = 'prisme.notifications';

function stored(): boolean {
	try {
		return localStorage.getItem(STORAGE_KEY) === 'on';
	} catch {
		return false;
	}
}

function worthNotifying(client: MatrixClient, event: MatrixEvent, room: Room): boolean {
	if (event.getSender() === client.getUserId() || !isMessage(event) || event.isRedacted()) return false;
	if (isDirect(client, room)) return !isMuted(client, room.roomId);
	return !!client.getPushActionsForEvent(event)?.tweaks?.highlight;
}

class NotifyState {
	enabled = $state(false);
	supported = $state(false);
	permission = $state<NotificationPermission>('default');

	/** Reads the stored toggle and the browser's permission. */
	restore() {
		this.supported = typeof Notification !== 'undefined';
		if (!this.supported) return;
		this.permission = Notification.permission;
		this.enabled = stored() && this.permission === 'granted';
	}

	/** Turns browser notifications on (asking for permission) or off. */
	async setEnabled(on: boolean) {
		if (on && this.supported && Notification.permission !== 'granted') {
			this.permission = await Notification.requestPermission();
		}
		this.enabled = on && this.permission === 'granted';
		try {
			localStorage.setItem(STORAGE_KEY, this.enabled ? 'on' : 'off');
		} catch {
			return;
		}
	}

	/** Shows a notification for live mentions and DMs while the tab is hidden. */
	watch(client: MatrixClient, open: (roomId: string) => void): () => void {
		const onTimeline = (event: MatrixEvent, room: Room | undefined, start?: boolean, removed?: boolean, data?: { liveEvent?: boolean }) => {
			if (!this.enabled || !document.hidden || !room || start || removed || !data?.liveEvent) return;
			if (!worthNotifying(client, event, room)) return;
			const sender = room.getMember(event.getSender() ?? '')?.name ?? event.getSender();
			const note = new Notification(`${sender} in ${room.name}`, { body: String(event.getContent().body ?? ''), tag: room.roomId });
			note.onclick = () => {
				window.focus();
				open(room.roomId);
			};
		};
		client.on(RoomEvent.Timeline, onTimeline);
		return () => client.off(RoomEvent.Timeline, onTimeline);
	}
}

export const notify = new NotifyState();
