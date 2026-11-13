import { Preferences } from '@capacitor/preferences';
import { invoke } from '@tauri-apps/api/core';
import { SecureStoragePlugin } from 'capacitor-secure-storage-plugin';
import { isApp, isDesktop } from './platform';

const TOKEN_KEY = 'prisme.token';

let cachedToken: string | null | undefined;

/**
 * Reads the bearer token from the platform's secure store (Android Keystore, macOS Keychain,
 * Windows Credential Manager, Secret Service), or null on the web. A token an older build kept in
 * Preferences moves to the secure store on the first read, so updating signs nobody out.
 */
export async function getToken(): Promise<string | null> {
	if (!isApp()) return null;
	if (cachedToken === undefined) cachedToken = (await readSecure()) ?? (await adoptPreferencesToken());
	return cachedToken;
}

/** Keeps `token` in the secure store. It stays in memory for this run even if the store refuses it. */
export async function setToken(token: string): Promise<void> {
	if (!isApp()) return;
	cachedToken = token;
	await writeSecure(token);
	await Preferences.remove({ key: TOKEN_KEY });
}

/** Forgets the bearer token, signing the installed app out. */
export async function clearToken(): Promise<void> {
	if (!isApp()) return;
	cachedToken = null;
	try {
		if (isDesktop()) await invoke('token_clear');
		else await SecureStoragePlugin.remove({ key: TOKEN_KEY });
	} catch (error) {
		console.warn('prisme: the secure store kept the token', error);
	}
	await Preferences.remove({ key: TOKEN_KEY });
}

async function readSecure(): Promise<string | null> {
	try {
		if (isDesktop()) return (await invoke<string | null>('token_get')) ?? null;
		return (await SecureStoragePlugin.get({ key: TOKEN_KEY })).value || null;
	} catch {
		return null;
	}
}

async function writeSecure(token: string): Promise<boolean> {
	try {
		if (isDesktop()) await invoke('token_set', { token });
		else await SecureStoragePlugin.set({ key: TOKEN_KEY, value: token });
		return true;
	} catch (error) {
		console.warn('prisme: the secure store refused the token', error);
		return false;
	}
}

/**
 * Moves the token builds before v0.20.0 kept in plain Preferences into the secure store. The
 * Preferences copy is deleted only once the secure store holds the token, so a device whose store
 * refuses it stays signed in as before.
 */
async function adoptPreferencesToken(): Promise<string | null> {
	const legacy = (await Preferences.get({ key: TOKEN_KEY })).value;
	if (!legacy) return null;
	if (await writeSecure(legacy)) await Preferences.remove({ key: TOKEN_KEY });
	return legacy;
}
