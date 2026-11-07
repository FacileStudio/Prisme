import { ExternalE2EEKeyProvider, Room, isE2EESupported } from 'livekit-client';
import { t } from '$lib/i18n/index.svelte';

/** Whether this browser can join an encrypted call at all. */
export function canEncrypt(): boolean {
	return typeof Worker !== 'undefined' && isE2EESupported();
}

/** Whether this browser can share a screen. Mobile browsers have no getDisplayMedia. */
export function canShareScreen(): boolean {
	return typeof navigator !== 'undefined' && typeof navigator.mediaDevices?.getDisplayMedia === 'function';
}

/**
 * Builds a LiveKit room that encrypts media and data with the room's key.
 * The worker is fresh for every call; the caller terminates it on leave.
 */
export async function encryptedRoom(key: string): Promise<{ room: Room; worker: Worker }> {
	const keyProvider = new ExternalE2EEKeyProvider();
	const worker = new Worker(new URL('livekit-client/e2ee-worker', import.meta.url), { type: 'module' });
	const room = new Room({ adaptiveStream: true, dynacast: true, encryption: { keyProvider, worker } });
	try {
		await keyProvider.setKey(key);
		await room.setE2EEEnabled(true);
	} catch (err) {
		worker.terminate();
		throw err;
	}
	return { room, worker };
}

/**
 * Turns a screen share failure into a message, or null when the person only
 * closed the picker. Chrome says "Permission denied by system" when the OS
 * blocks capture, and "Permission denied" for a cancel.
 */
export function shareError(err: unknown): string | null {
	const name = err instanceof Error ? err.name : '';
	const message = err instanceof Error ? err.message : '';
	if (name === 'NotAllowedError' && /system/i.test(message)) {
		return t('call.shareBlocked');
	}
	if (name === 'NotAllowedError' || name === 'AbortError') return null;
	return message || t('call.shareFailed');
}

/** How long a finished call lasted, the way the call card says it. */
export function callLength(ms: number): string {
	const minutes = Math.floor(ms / 60_000);
	if (minutes < 1) return t('call.underMinute');
	if (minutes < 60) return t('call.minutes', { count: minutes });
	const rest = minutes % 60;
	const hours = Math.floor(minutes / 60);
	return rest ? t('call.hoursMinutes', { hours, minutes: rest }) : t('call.hours', { count: hours });
}
