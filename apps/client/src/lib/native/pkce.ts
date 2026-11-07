const VERIFIER_ALPHABET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~';

export type AuthLink = { code: string; state: string; error: string };

/** Encodes bytes as unpadded base64url. */
export function base64Url(bytes: Uint8Array): string {
	let binary = '';
	for (const byte of bytes) binary += String.fromCharCode(byte);
	return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

/** Builds a PKCE verifier from the RFC 7636 unreserved characters. */
export function randomVerifier(length = 64): string {
	const bytes = crypto.getRandomValues(new Uint8Array(length));
	return Array.from(bytes, (b) => VERIFIER_ALPHABET[b % VERIFIER_ALPHABET.length]).join('');
}

/** Builds a login state from base64url, the only characters porte keeps on the app redirect. */
export function randomState(): string {
	return base64Url(crypto.getRandomValues(new Uint8Array(24)));
}

/** Returns the S256 PKCE challenge of a verifier. */
export async function pkceChallenge(verifier: string): Promise<string> {
	const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(verifier));
	return base64Url(new Uint8Array(digest));
}

/** Parses the `<scheme>://auth` sign-in callback, or returns null for any other link. */
export function parseAuthLink(raw: string, scheme: string): AuthLink | null {
	let url: URL;
	try {
		url = new URL(raw);
	} catch {
		return null;
	}
	if (url.protocol !== `${scheme}:` || url.hostname !== 'auth') return null;
	return {
		code: url.searchParams.get('code') ?? '',
		state: url.searchParams.get('state') ?? '',
		error: url.searchParams.get('error') ?? ''
	};
}
