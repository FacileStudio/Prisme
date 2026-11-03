import { describe, expect, it } from 'vitest';
import { countedLabel, railPages, type Nav, type SpaceNav } from './nav';

const quiet = { count: 0, strong: false, bold: false };

describe('rail pages', () => {
	it('labels unread rows with their count or a dot', () => {
		expect(countedLabel('general', quiet)).toBe('general');
		expect(countedLabel('general', { count: 0, strong: false, bold: true })).toBe('general ·');
		expect(countedLabel('general', { count: 120, strong: true, bold: true })).toBe('general · 99+');
	});

	it('lists browse, the current space channels, then DMs, and marks the open room', () => {
		const current = {
			name: 'Facile',
			joined: [{ roomId: '!a:x', name: 'general', badge: quiet, muted: false }]
		} as SpaceNav;
		const nav = {
			spaces: [current],
			dms: [{ roomId: '!d:x', name: 'Pixel', userId: '@pixel:x', invite: false, agent: true, badge: quiet }],
			dmBadge: quiet,
			invites: [{}]
		} as unknown as Nav;
		const pages = railPages(nav, current, '/room/%21a%3Ax', '!a:x');
		expect(pages.map((p) => [p.label, p.kind, p.group, p.active])).toEqual([
			['Overview · 1', 'browse', undefined, false],
			['general', 'channel', 'Facile', true],
			['Pixel', 'agent', 'Direct messages', false]
		]);
		expect(pages[1].href).toBe('/room/!a%3Ax');
	});
});
