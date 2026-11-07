import { RoomEvent, type MatrixClient, type MatrixEvent, type Room } from 'matrix-js-sdk';
import { t } from '$lib/i18n/index.svelte';
import { isDirect, isMuted } from '$lib/matrix/rooms';
import { isMessage } from '$lib/matrix/timeline';
import { desktopPermission, requestDesktop, showDesktop } from '$lib/native/desktop';
import { isDesktop, isNative } from '$lib/native/platform';
import { nativePermission, requestNative, startNativePush, stopNativePush } from '$lib/native/push';

const STORAGE_KEY = 'prisme.notifications';

/** Where notifications come from: the browser tab, the desktop app's sync, or Android push. */
export type Channel = 'web' | 'desktop' | 'native';

type Open = (roomId: string) => void;

function channel(): Channel {
	if (typeof window === 'undefined') return 'web';
	if (isNative()) return 'native';
	return isDesktop() ? 'desktop' : 'web';
}

function stored(): string | null {
	try {
		return localStorage.getItem(STORAGE_KEY);
	} catch {
		return null;
	}
}

function store(on: boolean) {
	try {
		localStorage.setItem(STORAGE_KEY, on ? 'on' : 'off');
	} catch {
		return;
	}
}

function worthNotifying(client: MatrixClient, event: MatrixEvent, room: Room): boolean {
	if (event.getSender() === client.getUserId() || !isMessage(event) || event.isRedacted()) return false;
	if (isDirect(client, room)) return !isMuted(client, room.roomId);
	return !!client.getPushActionsForEvent(event)?.tweaks?.highlight;
}

class NotifyState {
	readonly channel = channel();
	enabled = $state(false);
	supported = $state(false);
	available = $state(true);
	permission = $state<NotificationPermission>('default');
	private ready: Promise<void> = Promise.resolve();
	private opener: Open = () => undefined;

	/** Reads the stored toggle and the permission of this browser or app. */
	restore() {
		if (this.channel === 'web') {
			this.supported = typeof Notification !== 'undefined';
			if (!this.supported) return;
		}
		this.supported = true;
		this.ready = this.refresh();
	}

	private async refresh() {
		this.permission = await this.read().catch((): NotificationPermission => 'default');
		this.enabled = this.wanted() && this.permission === 'granted';
	}

	private wanted(): boolean {
		return this.channel === 'desktop' ? stored() !== 'off' : stored() === 'on';
	}

	private read(): Promise<NotificationPermission> {
		if (this.channel === 'native') return nativePermission();
		if (this.channel === 'desktop') return desktopPermission();
		return Promise.resolve(Notification.permission);
	}

	private ask(): Promise<NotificationPermission> {
		if (this.channel === 'native') return requestNative();
		if (this.channel === 'desktop') return requestDesktop();
		return Notification.requestPermission();
	}

	/** Turns notifications on (asking for permission) or off for this device. */
	async setEnabled(on: boolean) {
		if (on && this.supported && this.permission !== 'granted') {
			this.permission = await this.ask().catch((): NotificationPermission => 'denied');
		}
		this.enabled = on && this.permission === 'granted';
		store(this.enabled);
		if (this.channel !== 'native') return;
		if (this.enabled) this.available = await startNativePush(this.opener);
		else await stopNativePush();
	}

	/**
	 * Notifies of live mentions and direct messages. The web tab and the desktop app do it from
	 * their own sync while they are out of sight; the Android app hands that job to push.
	 */
	watch(client: MatrixClient, open: Open): () => void {
		this.opener = open;
		if (this.channel === 'native') {
			void this.ready.then(async () => {
				if (this.enabled) this.available = await startNativePush(open);
			});
			return () => undefined;
		}
		const onTimeline = (event: MatrixEvent, room: Room | undefined, start?: boolean, removed?: boolean, data?: { liveEvent?: boolean }) => {
			if (!this.enabled || !this.away() || !room || start || removed || !data?.liveEvent) return;
			if (!worthNotifying(client, event, room)) return;
			const sender = room.getMember(event.getSender() ?? '')?.rawDisplayName ?? event.getSender();
			this.show(t('notify.title', { sender: String(sender), room: room.name }), String(event.getContent().body ?? ''), room.roomId);
		};
		client.on(RoomEvent.Timeline, onTimeline);
		return () => client.off(RoomEvent.Timeline, onTimeline);
	}

	private away(): boolean {
		return this.channel === 'desktop' ? !document.hasFocus() : document.hidden;
	}

	private show(title: string, body: string, roomId: string) {
		if (this.channel === 'desktop') {
			showDesktop(title, body);
			return;
		}
		const note = new Notification(title, { body, tag: roomId });
		note.onclick = () => {
			window.focus();
			this.opener(roomId);
		};
	}
}

export const notify = new NotifyState();
