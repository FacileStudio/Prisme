import { describe, expect, it } from 'vitest';
import { parseAuthLink, pkceChallenge, randomState, randomVerifier } from './pkce';

describe('app sign-in helpers', () => {
	it('computes the RFC 7636 S256 challenge', async () => {
		expect(await pkceChallenge('dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk')).toBe('E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM');
	});

	it('builds states porte keeps and verifiers within the RFC length', () => {
		expect(randomState()).toMatch(/^[A-Za-z0-9_-]{32}$/);
		expect(randomVerifier()).toMatch(/^[A-Za-z0-9._~-]{64}$/);
	});

	it('reads only the auth callback of its own scheme', () => {
		expect(parseAuthLink('studio.facile.prisme://auth?code=c&state=s', 'studio.facile.prisme')).toEqual({ code: 'c', state: 's', error: '' });
		expect(parseAuthLink('other.app://auth?code=c', 'studio.facile.prisme')).toBeNull();
		expect(parseAuthLink('studio.facile.prisme://room/x', 'studio.facile.prisme')).toBeNull();
	});
});
