import { beforeEach, describe, expect, it, vi } from 'vitest';

const device = vi.hoisted(() => ({
	native: true,
	desktop: false,
	refuse: false,
	preferences: new Map<string, string>(),
	secure: new Map<string, string>(),
	keychain: null as string | null
}));

vi.mock('@capacitor/core', () => ({ Capacitor: { isNativePlatform: () => device.native } }));

vi.mock('@capacitor/preferences', () => ({
	Preferences: {
		get: async ({ key }: { key: string }) => ({ value: device.preferences.get(key) ?? null }),
		remove: async ({ key }: { key: string }) => void device.preferences.delete(key)
	}
}));

vi.mock('capacitor-secure-storage-plugin', () => ({
	SecureStoragePlugin: {
		get: async ({ key }: { key: string }) => {
			const value = device.secure.get(key);
			if (value === undefined) throw new Error('Item with given key does not exist');
			return { value };
		},
		set: async ({ key, value }: { key: string; value: string }) => {
			if (device.refuse) throw new Error('keystore unavailable');
			device.secure.set(key, value);
			return { value: true };
		},
		remove: async ({ key }: { key: string }) => ({ value: device.secure.delete(key) })
	}
}));

vi.mock('@tauri-apps/api/core', () => ({
	isTauri: () => device.desktop,
	invoke: async (command: string, args?: { token: string }) => {
		if (command === 'token_get') return device.keychain;
		if (device.refuse) throw new Error('no secret service');
		device.keychain = command === 'token_set' ? (args?.token ?? null) : null;
		return null;
	}
}));

async function platform() {
	vi.resetModules();
	return import('./token');
}

beforeEach(() => {
	Object.assign(device, { native: true, desktop: false, refuse: false, keychain: null });
	device.preferences.clear();
	device.secure.clear();
	vi.spyOn(console, 'warn').mockImplementation(() => undefined);
});

describe('the app token', () => {
	it('moves an Android token from Preferences to the secure store on the first read', async () => {
		device.preferences.set('prisme.token', 'old-token');
		const { getToken } = await platform();
		expect(await getToken()).toBe('old-token');
		expect(device.secure.get('prisme.token')).toBe('old-token');
		expect(device.preferences.has('prisme.token')).toBe(false);
		const again = await platform();
		expect(await again.getToken()).toBe('old-token');
	});

	it('moves a desktop token from Preferences to the keychain', async () => {
		Object.assign(device, { native: false, desktop: true });
		device.preferences.set('prisme.token', 'desk-token');
		const { getToken } = await platform();
		expect(await getToken()).toBe('desk-token');
		expect(device.keychain).toBe('desk-token');
		expect(device.preferences.has('prisme.token')).toBe(false);
	});

	it('keeps the Preferences token when the secure store refuses it', async () => {
		device.refuse = true;
		device.preferences.set('prisme.token', 'old-token');
		const { getToken } = await platform();
		expect(await getToken()).toBe('old-token');
		expect(device.preferences.get('prisme.token')).toBe('old-token');
		expect(device.secure.size).toBe(0);
	});

	it('stores a new token only in the secure store and forgets it on sign-out', async () => {
		device.preferences.set('prisme.token', 'stale');
		const { setToken, clearToken, getToken } = await platform();
		await setToken('new-token');
		expect(device.secure.get('prisme.token')).toBe('new-token');
		expect(device.preferences.has('prisme.token')).toBe(false);
		await clearToken();
		expect(device.secure.size).toBe(0);
		expect(await getToken()).toBeNull();
	});

	it('leaves the web build alone', async () => {
		Object.assign(device, { native: false, desktop: false });
		device.preferences.set('prisme.token', 'web');
		const { getToken, setToken, clearToken } = await platform();
		expect(await getToken()).toBeNull();
		await setToken('t');
		await clearToken();
		expect(device.preferences.get('prisme.token')).toBe('web');
		expect(device.secure.size).toBe(0);
	});
});
