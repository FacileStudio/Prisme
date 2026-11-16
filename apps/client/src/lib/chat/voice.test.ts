import { describe, expect, it } from 'vitest';
import { bars, clock, downsample, extensionOf, recordingType, voiceContent, voiceOf, WAVEFORM_MAX } from './voice';

describe('voice waveform', () => {
	it('keeps the peak of each bucket and normalises to the loudest', () => {
		const out = downsample([0.1, 0.5, 0.2, 0.25, 0, 0, 0.1, 0.05], 4);
		expect(out).toEqual([WAVEFORM_MAX, 512, 0, 204]);
	});

	it('draws flat bars for silence and an empty recording', () => {
		expect(downsample([], 3)).toEqual([0, 0, 0]);
		expect(downsample([0, 0, 0, 0], 2)).toEqual([0, 0]);
		expect(bars([], 2)).toEqual([0.08, 0.08]);
	});

	it('stretches fewer samples than bars', () => {
		expect(downsample([1, 0.5], 4)).toEqual([1024, 1024, 512, 512]);
		expect(bars([1024, 512], 4)).toEqual([1, 1, 0.5, 0.5]);
	});
});

describe('voice content', () => {
	it('round trips a voice message', () => {
		const content = voiceContent('mxc://s/id', { type: 'audio/webm;codecs=opus', size: 2048 }, 4200, [1, 2, 3], 'Voice message');
		const voice = voiceOf(content);
		expect(voice).toEqual({ url: 'mxc://s/id', duration: 4200, waveform: [1, 2, 3], mimetype: 'audio/webm', size: 2048 });
	});

	it('leaves plain audio files to the file card', () => {
		expect(voiceOf({ msgtype: 'm.audio', url: 'mxc://s/id', body: 'song.mp3' })).toBeNull();
		expect(voiceOf({ msgtype: 'm.text', body: 'hi' })).toBeNull();
	});

	it('formats clocks and picks formats', () => {
		expect(clock(7000)).toBe('0:07');
		expect(clock(102_400)).toBe('1:42');
		expect(recordingType((t) => t === 'audio/mp4')).toBe('audio/mp4');
		expect(recordingType(() => false)).toBe('');
		expect(extensionOf('audio/mp4')).toBe('m4a');
		expect(extensionOf('audio/webm')).toBe('webm');
	});
});
