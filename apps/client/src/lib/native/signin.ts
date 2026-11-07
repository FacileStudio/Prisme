import { App } from '@capacitor/app';
import { Browser } from '@capacitor/browser';
import { getCurrentWindow } from '@tauri-apps/api/window';
import { onOpenUrl } from '@tauri-apps/plugin-deep-link';
import { openUrl } from '@tauri-apps/plugin-opener';
import { goto } from '$app/navigation';
import { t } from '$lib/i18n/index.svelte';
import { parseAuthLink, pkceChallenge, randomState, randomVerifier } from './pkce';
import { APP_SCHEME, apiBase, isDesktop, setToken } from './platform';

const PENDING_KEY = 'prisme.pending-login';

type Pending = { verifier: string; state: string };

let registered = false;

function remember(value: Pending) {
	try {
		sessionStorage.setItem(PENDING_KEY, JSON.stringify(value));
	} catch {
		return;
	}
}

function recall(): Pending | null {
	try {
		const raw = sessionStorage.getItem(PENDING_KEY);
		sessionStorage.removeItem(PENDING_KEY);
		return raw ? (JSON.parse(raw) as Pending) : null;
	} catch {
		return null;
	}
}

/** Opens porte's app sign-in in the system browser, bound to a fresh PKCE verifier. */
export async function startAppLogin(): Promise<void> {
	const verifier = randomVerifier();
	const state = randomState();
	remember({ verifier, state });
	const query = new URLSearchParams({ flow: 'mobile', challenge: await pkceChallenge(verifier), state });
	const url = `${apiBase()}/api/auth/oidc?${query}`;
	if (isDesktop()) await openUrl(url);
	else await Browser.open({ url });
}

async function exchange(code: string, state: string): Promise<void> {
	const saved = recall();
	if (!saved || saved.state !== state) throw new Error(t('app.signInMismatch'));
	const response = await fetch(`${apiBase()}/api/auth/oidc/exchange`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ code, verifier: saved.verifier })
	});
	const body = (await response.json().catch(() => ({}))) as { token?: string };
	if (!response.ok || !body.token) throw new Error(t('app.signInFailed'));
	await setToken(body.token);
}

async function finish(code: string, state: string, error: string) {
	let failure = error;
	if (!failure) failure = await exchange(code, state).then(() => '', (err: unknown) => (err instanceof Error ? err.message : t('app.signInFailed')));
	if (isDesktop()) await getCurrentWindow().setFocus().catch(() => undefined);
	else await Browser.close().catch(() => undefined);
	if (failure) await goto(`/login?error=${encodeURIComponent(failure)}`, { replaceState: true });
	else await goto('/welcome', { replaceState: true });
}

function handle(url: string) {
	const link = parseAuthLink(url, APP_SCHEME);
	if (link) void finish(link.code, link.state, link.error);
}

/** Listens for the sign-in callback once per app start. */
export function registerDeepLinks(): void {
	if (registered) return;
	registered = true;
	if (isDesktop()) void onOpenUrl((urls) => urls.forEach(handle));
	else void App.addListener('appUrlOpen', ({ url }) => handle(url));
}
