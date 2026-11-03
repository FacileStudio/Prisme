import { describe, expect, it } from 'vitest';
import { callLength, shareError } from './media';

describe('screen share errors', () => {
	it('stays silent on a closed picker and explains real failures', () => {
		expect(shareError(new DOMException('Permission denied', 'NotAllowedError'))).toBeNull();
		expect(shareError(new DOMException('The user aborted a request.', 'AbortError'))).toBeNull();
		expect(shareError(new DOMException('Permission denied by system', 'NotAllowedError'))).toMatch(/system settings/);
		expect(shareError(new DOMException('Could not start video source', 'NotReadableError'))).toBe(
			'Could not start video source'
		);
	});
});

describe('call length', () => {
	it('reads like a person would say it', () => {
		expect(callLength(30_000)).toBe('Under a minute');
		expect(callLength(12 * 60_000)).toBe('12 min');
		expect(callLength(60 * 60_000)).toBe('1 h');
		expect(callLength(75 * 60_000)).toBe('1 h 15 min');
	});
});
