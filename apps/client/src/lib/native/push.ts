import { PushNotifications } from '@capacitor/push-notifications';
import { prisme } from '$lib/api/prisme';

let listening = false;
let opener: (roomId: string) => void = () => undefined;

/** Reads whether Android lets Prisme post notifications. */
export async function nativePermission(): Promise<NotificationPermission> {
	const { receive } = await PushNotifications.checkPermissions();
	return receive === 'granted' ? 'granted' : receive === 'denied' ? 'denied' : 'default';
}

/** Asks Android for the notification permission. */
export async function requestNative(): Promise<NotificationPermission> {
	const { receive } = await PushNotifications.requestPermissions();
	return receive === 'granted' ? 'granted' : receive === 'denied' ? 'denied' : 'default';
}

function listen() {
	if (listening) return;
	listening = true;
	void PushNotifications.addListener('registration', ({ value }) => {
		void prisme.registerPushDevice(value).catch(() => undefined);
	});
	void PushNotifications.addListener('pushNotificationActionPerformed', ({ notification }) => {
		const roomId = notification.data?.room_id;
		if (typeof roomId === 'string' && roomId) opener(roomId);
	});
}

/**
 * Registers this phone with Firebase and hands the token to Prisme, which points a Synapse
 * pusher at it. It never registers while the server has no push configured, so the app only
 * talks to Firebase when something will use it. Returns whether push is on.
 */
export async function startNativePush(open: (roomId: string) => void): Promise<boolean> {
	opener = open;
	const config = await prisme.pushConfig().catch(() => null);
	if (!config?.native) return false;
	listen();
	await PushNotifications.register();
	return true;
}

/** Drops the Firebase token, so the next push fails and Synapse removes the pusher. */
export async function stopNativePush(): Promise<void> {
	await PushNotifications.unregister().catch(() => undefined);
}
