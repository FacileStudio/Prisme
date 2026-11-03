import { describe, expect, it } from 'vitest';
import { shareError } from './media';

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
