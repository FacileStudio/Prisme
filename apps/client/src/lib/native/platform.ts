import { Capacitor } from '@capacitor/core';
import { Preferences } from '@capacitor/preferences';
import { isTauri } from '@tauri-apps/api/core';

export const APP_SCHEME = 'studio.facile.prisme';

const TOKEN_KEY = 'prisme.token';
const APP_API_ORIGIN = 'https://prisme.facile.studio';

let cachedToken: string | null | undefined;

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

/** Reads the bearer token an installed app keeps, or null on the web. */
export async function getToken(): Promise<string | null> {
	if (!isApp()) return null;
	if (cachedToken === undefined) cachedToken = (await Preferences.get({ key: TOKEN_KEY })).value;
	return cachedToken;
}

/** Stores the bearer token the app sign-in returned. */
export async function setToken(token: string): Promise<void> {
	cachedToken = token;
	await Preferences.set({ key: TOKEN_KEY, value: token });
}

/** Forgets the bearer token, signing the installed app out. */
export async function clearToken(): Promise<void> {
	if (!isApp()) return;
	cachedToken = null;
	await Preferences.remove({ key: TOKEN_KEY });
}
