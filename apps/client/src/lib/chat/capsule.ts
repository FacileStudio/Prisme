import { prisme } from '$lib/api/prisme';
import { fromBase64, toBase64 } from './base64';

export const CAPSULE_KEY = 'studio.facile.capsule';
export const EXPIRIES = ['1h', '24h', '7d', '30d'] as const;
export type Expiry = (typeof EXPIRIES)[number];

/** The studio.facile.capsule block of a sealed message. `url` carries the key in its fragment. */
export type CapsuleInfo = { id: string; url: string; expires_at?: string };

/** Capsule's plaintext envelope for a text-only capsule. */
export function envelope(text: string): string {
	return JSON.stringify({ __capsule: 1, text });
}

/** The text of a decrypted Capsule payload, enveloped or a bare legacy string. */
export function envelopeText(decrypted: string): string {
	if (!decrypted.startsWith('{"__capsule":1')) return decrypted;
	try {
		const parsed = JSON.parse(decrypted) as { __capsule?: number; text?: unknown };
		return parsed.__capsule === 1 && typeof parsed.text === 'string' ? parsed.text : decrypted;
	} catch {
		return decrypted;
	}
}

/** Encrypts text with a fresh AES-256-GCM key: base64(iv || ciphertext), and the base64url key. */
export async function encryptText(text: string): Promise<{ ciphertext: string; key: string }> {
	const key = await crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, ['encrypt', 'decrypt']);
	const iv = crypto.getRandomValues(new Uint8Array(12));
	const sealed = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, new TextEncoder().encode(envelope(text)));
	const combined = new Uint8Array(iv.length + sealed.byteLength);
	combined.set(iv, 0);
	combined.set(new Uint8Array(sealed), iv.length);
	const raw = new Uint8Array(await crypto.subtle.exportKey('raw', key));
	return { ciphertext: toBase64(combined), key: toBase64(raw, true) };
}

/** Decrypts a Capsule ciphertext with its base64url key and returns the text. */
export async function decryptText(ciphertext: string, key: string): Promise<string> {
	const imported = await crypto.subtle.importKey('raw', fromBase64(key), { name: 'AES-GCM', length: 256 }, false, ['decrypt']);
	const combined = fromBase64(ciphertext);
	const plain = await crypto.subtle.decrypt({ name: 'AES-GCM', iv: combined.slice(0, 12) }, imported, combined.slice(12));
	return envelopeText(new TextDecoder().decode(plain));
}

/** The capsule block of raw message content, when it names a capsule with its key. */
export function capsuleOf(content: Record<string, unknown>): CapsuleInfo | null {
	const block = content[CAPSULE_KEY] as CapsuleInfo | undefined;
	return typeof block?.id === 'string' && typeof block.url === 'string' ? block : null;
}

/** The key in a capsule link's fragment, or an empty string. */
export function keyOf(url: string): string {
	const at = url.indexOf('#');
	return at === -1 ? '' : url.slice(at + 1);
}

/** Seals text in Capsule through Prisme and returns the message block. */
export async function sealText(text: string, expiresIn: Expiry): Promise<CapsuleInfo> {
	const { ciphertext, key } = await encryptText(text);
	const sealed = await prisme.sealCapsule(ciphertext, expiresIn);
	return { id: sealed.id, url: `${sealed.url}#${key}`, expires_at: sealed.expires_at };
}

/** Opens a capsule, which burns it, and returns its text. */
export async function openCapsule(info: CapsuleInfo): Promise<string> {
	const { ciphertext } = await prisme.openCapsule(info.id);
	return decryptText(ciphertext, keyOf(info.url));
}
