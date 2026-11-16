import { describe, expect, it } from 'vitest';
import { fileSize, nuageBlock, nuageBody, nuageOf, NUAGE_KEY } from './nuage';

const shared = { file_id: 7, name: 'plan.pdf', mime_type: 'application/pdf', size: 2048, url: 'https://nuage.facile.studio/s/AbC_12-x' };

describe('nuageBlock', () => {
	it('keeps a Nuage share page and renames the mime field', () => {
		expect(nuageBlock(shared)).toEqual({ file_id: 7, name: 'plan.pdf', size: 2048, mime: 'application/pdf', url: shared.url });
	});

	it('refuses any link that is not a plain Nuage share page', () => {
		for (const url of [
			'https://evil.example/s/AbC',
			'http://nuage.facile.studio/s/AbC',
			'https://nuage.facile.studio/files/7',
			'https://nuage.facile.studio/s/AbC?x=1',
			'https://nuage.facile.studio/s/AbC#k',
			'not a url'
		]) {
			expect(nuageBlock({ ...shared, url })).toBeNull();
		}
	});
});

describe('nuage content', () => {
	it('reads back the block a message carries and gives other clients a readable body', () => {
		const info = nuageBlock(shared)!;
		expect(nuageOf({ msgtype: 'm.text', body: nuageBody(info), [NUAGE_KEY]: info })).toEqual(info);
		expect(nuageBody(info)).toBe('plan.pdf https://nuage.facile.studio/s/AbC_12-x');
		expect(nuageOf({ [NUAGE_KEY]: { ...info, url: 'javascript:alert(1)' } })).toBeNull();
		expect(nuageOf({ body: 'hi' })).toBeNull();
	});
});

describe('fileSize', () => {
	it('reads like people write sizes', () => {
		expect(fileSize(0)).toBe('0 B');
		expect(fileSize(2048)).toBe('2 KB');
		expect(fileSize(3.4 * 1024 * 1024)).toBe('3.4 MB');
		expect(fileSize(15 * 1024 * 1024 * 1024)).toBe('15 GB');
	});
});
