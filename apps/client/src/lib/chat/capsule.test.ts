import { describe, expect, it } from 'vitest';
import { decryptText, encryptText, envelope, envelopeText, keyOf } from './capsule';
import { fromBase64 } from './base64';

describe('capsule', () => {
	it('seals and opens a round trip in Capsule format', async () => {
		const { ciphertext, key } = await encryptText('the door code is 4521 ✓');
		expect(key).not.toMatch(/[+/=]/);
		expect(fromBase64(key).length).toBe(32);
		expect(fromBase64(ciphertext).length).toBeGreaterThan(12 + 16);
		expect(await decryptText(ciphertext, key)).toBe('the door code is 4521 ✓');
	});

	it('refuses the wrong key', async () => {
		const { ciphertext } = await encryptText('secret');
		const other = await encryptText('x');
		await expect(decryptText(ciphertext, other.key)).rejects.toThrow();
	});

	it('reads the envelope and legacy strings', () => {
		expect(envelope('hi')).toBe('{"__capsule":1,"text":"hi"}');
		expect(envelopeText('{"__capsule":1,"text":"hi"}')).toBe('hi');
		expect(envelopeText('plain')).toBe('plain');
	});

	it('takes the key from the fragment only', () => {
		expect(keyOf('https://capsule.facile.studio/cap_0123456789abcdef#abc')).toBe('abc');
		expect(keyOf('https://capsule.facile.studio/cap_0123456789abcdef')).toBe('');
	});
});
