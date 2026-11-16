export const WAVEFORM_BARS = 48;
export const WAVEFORM_MAX = 1024;
export const MAX_VOICE_MS = 10 * 60 * 1000;
export const MIN_VOICE_MS = 700;

export type VoiceInfo = { url: string; duration: number; waveform: number[]; mimetype: string; size: number };

const VOICE_KEYS = ['org.matrix.msc3245.voice', 'm.voice'];
const AUDIO_KEYS = ['org.matrix.msc1767.audio', 'm.audio'];
const TYPES = ['audio/webm;codecs=opus', 'audio/ogg;codecs=opus', 'audio/mp4', 'audio/webm'];

/**
 * Shrinks a run of loudness samples (0 to 1) into `bars` peaks scaled to the
 * MSC3246 range of 0 to 1024, so a long recording and a short one draw the same.
 */
export function downsample(samples: number[], bars = WAVEFORM_BARS): number[] {
	if (samples.length === 0) return Array.from({ length: bars }, () => 0);
	const out: number[] = [];
	for (let i = 0; i < bars; i++) {
		const from = Math.floor((i * samples.length) / bars);
		const to = Math.max(from + 1, Math.floor(((i + 1) * samples.length) / bars));
		let peak = 0;
		for (let j = from; j < to && j < samples.length; j++) peak = Math.max(peak, samples[j]);
		out.push(Math.round(Math.min(1, Math.max(0, peak)) * WAVEFORM_MAX));
	}
	const top = Math.max(...out);
	return top === 0 ? out : out.map((v) => Math.round((v / top) * WAVEFORM_MAX));
}

/** Resamples a stored waveform to the number of bars a player draws, as heights from 0 to 1. */
export function bars(waveform: number[], count = WAVEFORM_BARS): number[] {
	if (waveform.length === 0) return Array.from({ length: count }, () => 0.08);
	return Array.from({ length: count }, (_, i) => {
		const v = waveform[Math.min(waveform.length - 1, Math.floor((i * waveform.length) / count))] ?? 0;
		return Math.max(0.08, Math.min(1, v / WAVEFORM_MAX));
	});
}

/** A clock reading for a duration in milliseconds: 0:07, 1:42, 12:05. */
export function clock(ms: number): string {
	const total = Math.max(0, Math.round(ms / 1000));
	return `${Math.floor(total / 60)}:${String(total % 60).padStart(2, '0')}`;
}

/** The voice message an `m.audio` content carries, or null for any other content. */
export function voiceOf(content: Record<string, unknown>): VoiceInfo | null {
	if (content.msgtype !== 'm.audio' || typeof content.url !== 'string') return null;
	if (!VOICE_KEYS.some((k) => k in content)) return null;
	const audio = (AUDIO_KEYS.map((k) => content[k]).find((v) => v && typeof v === 'object') ?? {}) as {
		duration?: unknown;
		waveform?: unknown;
	};
	const info = (content.info ?? {}) as { duration?: unknown; mimetype?: unknown; size?: unknown };
	const duration = Number(audio.duration ?? info.duration) || 0;
	const waveform = Array.isArray(audio.waveform) ? audio.waveform.map(Number).filter(Number.isFinite) : [];
	return {
		url: content.url,
		duration,
		waveform,
		mimetype: typeof info.mimetype === 'string' ? info.mimetype : 'audio/webm',
		size: Number(info.size) || 0
	};
}

/** The `m.audio` content of a voice message, readable by Element and other MSC3245 clients. */
export function voiceContent(url: string, blob: { type: string; size: number }, duration: number, waveform: number[], body: string) {
	const mimetype = blob.type.split(';')[0] || 'audio/webm';
	return {
		msgtype: 'm.audio',
		body,
		url,
		info: { mimetype, size: blob.size, duration },
		'org.matrix.msc1767.audio': { duration, waveform },
		'org.matrix.msc1767.text': body,
		'org.matrix.msc3245.voice': {}
	};
}

/** The first recording format this browser supports, or an empty string to let it choose. */
export function recordingType(supports: (type: string) => boolean): string {
	return TYPES.find((type) => supports(type)) ?? '';
}

/** The file extension that matches a recording's mimetype. */
export function extensionOf(mimetype: string): string {
	if (mimetype.includes('mp4')) return 'm4a';
	if (mimetype.includes('ogg')) return 'ogg';
	return 'webm';
}
