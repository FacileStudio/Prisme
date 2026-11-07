import { describe, expect, it } from 'vitest';
import { detectPlatform, ordered } from './downloads';

describe('downloads', () => {
	it('tells Android from desktop Linux and leaves iOS out', () => {
		expect(detectPlatform('Mozilla/5.0 (Linux; Android 14; Pixel 8)')).toBe('android');
		expect(detectPlatform('Mozilla/5.0 (X11; Linux x86_64)')).toBe('linux');
		expect(detectPlatform('Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)')).toBeNull();
		expect(detectPlatform('Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0)')).toBe('macos');
	});

	it('lists the visitor platform first', () => {
		expect(ordered('linux').map((b) => b.platform)).toEqual(['linux', 'android', 'macos', 'windows']);
		expect(ordered(null).map((b) => b.platform)).toEqual(['android', 'macos', 'windows', 'linux']);
	});
});
