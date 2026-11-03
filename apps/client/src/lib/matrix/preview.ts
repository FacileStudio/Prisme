import type { Room } from 'matrix-js-sdk';
import { stripReplyFallback } from '$lib/chat/markdown';
import { isMessage, type CallInfo } from './timeline';

export type Preview = { text: string; at: number };

/** The one line a conversation list shows for a message's content. */
export function previewText(content: Record<string, unknown>, mine: boolean): string {
	const call = content['studio.facile.call'] as CallInfo | undefined;
	const kinds: Record<string, string> = { 'm.image': 'Sent an image', 'm.file': 'Sent a file', 'm.video': 'Sent a video', 'm.audio': 'Sent audio' };
	const text = call?.call_id
		? call.ended_at ? 'Call ended' : 'Started a call'
		: (kinds[String(content.msgtype)] ?? stripReplyFallback(String(content.body ?? '')).replace(/\s+/g, ' ').trim());
	if (!text) return 'Deleted message';
	return mine ? `You: ${text}` : text;
}

/** The newest message of a room's live timeline, as a list preview. */
export function lastMessage(room: Room, me: string): Preview | null {
	const events = room.getLiveTimeline().getEvents();
	for (let index = events.length - 1; index >= 0; index--) {
		const event = events[index];
		if (!isMessage(event) || event.isRedacted()) continue;
		return { text: previewText(event.getContent(), event.getSender() === me), at: event.getTs() };
	}
	return null;
}

/** A message time for a list: the hour today, the weekday this week, the date before. */
export function shortTime(ts: number, now = Date.now()): string {
	const date = new Date(ts);
	const today = new Date(now);
	if (date.toDateString() === today.toDateString()) return date.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
	if (now - ts < 6 * 86_400_000) return date.toLocaleDateString(undefined, { weekday: 'short' });
	const year = date.getFullYear() === today.getFullYear() ? undefined : 'numeric';
	return date.toLocaleDateString(undefined, { day: 'numeric', month: 'short', year });
}
