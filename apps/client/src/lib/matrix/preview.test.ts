import { describe, expect, it } from 'vitest';
import { previewText, shortTime } from './preview';

describe('conversation preview', () => {
	it('shows one clean line per kind of message', () => {
		expect(previewText({ msgtype: 'm.text', body: '> <@a:x> old\n\nhello\nthere' }, false)).toBe('hello there');
		expect(previewText({ msgtype: 'm.text', body: 'on my way' }, true)).toBe('You: on my way');
		expect(previewText({ msgtype: 'm.image', body: 'cat.png' }, false)).toBe('Sent an image');
		expect(previewText({ msgtype: 'm.text', body: 'x', 'studio.facile.call': { call_id: 'c', started_by: '@a:x' } }, false)).toBe('Started a call');
		expect(previewText({ msgtype: 'm.text', body: 'x', 'studio.facile.call': { call_id: 'c', started_by: '@a:x', ended_at: 1 } }, true)).toBe('You: Call ended');
		expect(previewText({}, false)).toBe('Deleted message');
	});
});

describe('short time', () => {
	it('gets coarser as messages age', () => {
		const now = new Date(2026, 10, 3, 15, 0).getTime();
		expect(shortTime(new Date(2026, 10, 3, 9, 5).getTime(), now)).toMatch(/09.05/);
		expect(shortTime(new Date(2026, 10, 1, 9, 5).getTime(), now)).not.toMatch(/9/);
		expect(shortTime(new Date(2025, 0, 1).getTime(), now)).toMatch(/2025/);
	});
});
