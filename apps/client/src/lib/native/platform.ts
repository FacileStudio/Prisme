import { Capacitor } from '@capacitor/core';
import { isTauri } from '@tauri-apps/api/core';

export const APP_SCHEME = 'studio.facile.prisme';

const APP_API_ORIGIN = 'https://prisme.facile.studio';

/** True inside the Capacitor Android app. */
export function isNative(): boolean {
	return Capacitor.isNativePlatform();
}

/** True inside the Tauri desktop app. */
export function isDesktop(): boolean {
	return isTauri();
}

/** True on macOS in the desktop app, where the window keeps the system's traffic lights. */
export function hasNativeTrafficLights(): boolean {
	return isDesktop() && /mac/i.test(navigator.userAgent);
}

/** True in either installed app, which signs in with a bearer token instead of a cookie. */
export function isApp(): boolean {
	return isNative() || isDesktop();
}

/**
 * The origin API paths are appended to: empty on the web, which stays same-origin, and
 * production for the installed apps, which serve their pages from their own origin.
 */
export function apiBase(): string {
	const override = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '');
	if (override) return override;
	return isApp() ? APP_API_ORIGIN : '';
}
