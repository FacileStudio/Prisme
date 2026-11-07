import { toast } from '@facile/muse';
import { isPermissionGranted, requestPermission, sendNotification } from '@tauri-apps/plugin-notification';
import { openUrl } from '@tauri-apps/plugin-opener';
import { relaunch } from '@tauri-apps/plugin-process';
import { check, type Update } from '@tauri-apps/plugin-updater';
import { t } from '$lib/i18n/index.svelte';

/** Reads whether the desktop OS lets Prisme show notifications. */
export async function desktopPermission(): Promise<NotificationPermission> {
	return (await isPermissionGranted()) ? 'granted' : 'default';
}

/** Asks the desktop OS for permission to show notifications. */
export async function requestDesktop(): Promise<NotificationPermission> {
	if (await isPermissionGranted()) return 'granted';
	return requestPermission();
}

/** Shows a system notification from the desktop app. */
export function showDesktop(title: string, body: string): void {
	sendNotification({ title, body });
}

async function install(update: Update) {
	try {
		await update.downloadAndInstall();
		await relaunch();
	} catch {
		toast.danger(t('app.updateFailed'));
	}
}

async function checkForUpdate() {
	const update = await check().catch(() => null);
	if (!update) return;
	toast.info(t('app.updateReady', { version: update.version }), {
		duration: 0,
		action: { label: t('app.updateInstall'), onClick: () => void install(update) }
	});
}

function openExternally(event: MouseEvent) {
	const anchor = (event.target as Element | null)?.closest?.('a[href]');
	if (!(anchor instanceof HTMLAnchorElement)) return;
	const url = new URL(anchor.href, location.href);
	if (url.origin === location.origin || !/^https?:$/.test(url.protocol)) return;
	event.preventDefault();
	void openUrl(url.href);
}

/**
 * Sets the desktop window up: marks the page for the title bar's offset, opens outside links
 * in the system browser instead of the app's webview, and offers a pending update.
 */
export function startDesktop(): void {
	document.documentElement.dataset.desktop = '';
	document.addEventListener('click', openExternally);
	void checkForUpdate();
}
