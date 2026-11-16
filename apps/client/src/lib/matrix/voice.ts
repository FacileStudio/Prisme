import type { MatrixClient, MatrixEvent } from 'matrix-js-sdk';
import { extensionOf, voiceContent, voiceOf, type VoiceInfo } from '$lib/chat/voice';
import { relation, type Target } from './messages';

/** Uploads a recording and posts it as an MSC3245 voice message. */
export async function sendVoice(
	client: MatrixClient,
	target: Target,
	recording: { blob: Blob; duration: number; waveform: number[] },
	body: string,
	lastInThread?: string
) {
	const { blob, duration, waveform } = recording;
	const type = blob.type.split(';')[0] || 'audio/webm';
	const name = `voice-${new Date().toISOString().replace(/[:.]/g, '-')}.${extensionOf(type)}`;
	const upload = await client.uploadContent(blob, { type, name });
	const content: Record<string, unknown> = voiceContent(upload.content_uri, blob, duration, waveform, body);
	const rel = relation(target, lastInThread);
	if (rel) content['m.relates_to'] = rel;
	await client.sendMessage(target.roomId, target.threadId ?? null, content as never);
}

/** The voice message an event carries, if it is one. */
export function voiceInfo(event: MatrixEvent): VoiceInfo | null {
	return event.getType() === 'm.room.message' ? voiceOf(event.getContent()) : null;
}
