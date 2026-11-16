import { describe, expect, it, vi } from 'vitest';

vi.mock('$lib/api/prisme', () => ({ prisme: { unfurl: vi.fn() } }));

import { suiteLinks, withoutFragment } from './links';

describe('suiteLinks', () => {
	it('finds suite links in order, once each, and trims punctuation', () => {
		const text =
			'Sheet: https://feuilles.facile.studio/s/abc, call https://echo.facile.studio/room/team#key=s3cret. ' +
			'Again https://feuilles.facile.studio/s/abc and (https://diapo.facile.studio/p/tok)';
		expect(suiteLinks(text)).toEqual([
			{ url: 'https://feuilles.facile.studio/s/abc', app: 'feuilles' },
			{ url: 'https://echo.facile.studio/room/team#key=s3cret', app: 'echo' },
			{ url: 'https://diapo.facile.studio/p/tok', app: 'diapo' }
		]);
	});

	it('never lists capsule, foreign, insecure or ported links', () => {
		const text =
			'https://capsule.facile.studio/cap_1#k http://docs.facile.studio/d/1 https://docs.facile.studio:8443/d ' +
			'https://evil.example/https://docs.facile.studio https://docs.facile.studio.evil.example/x';
		expect(suiteLinks(text)).toEqual([]);
	});

	it('stops at the limit', () => {
		const text = ['a', 'b', 'c', 'd'].map((id) => `https://clap.facile.studio/v/${id}`).join(' ');
		expect(suiteLinks(text, 2)).toHaveLength(2);
	});

	it('drops the fragment before a link leaves the browser', () => {
		expect(withoutFragment('https://echo.facile.studio/room/team#key=s3cret')).toBe('https://echo.facile.studio/room/team');
	});
});
