import { EventType, MsgType, RelationType, THREAD_RELATION_TYPE, type MatrixClient, type MatrixEvent } from 'matrix-js-sdk';
import { composeMessage } from '$lib/chat/compose';
import type { Mention } from '$lib/chat/mentions';

export type Draft = { text: string; mentions: Mention[] };
export type Target = { roomId: string; threadId?: string | null; replyTo?: MatrixEvent | null };

type Content = Record<string, unknown>;

function textContent(draft: Draft): Content {
	const { body, html, userIds } = composeMessage(draft.text, draft.mentions);
	const content: Content = { msgtype: MsgType.Text, body, 'm.mentions': { user_ids: userIds } };
	if (html !== null) Object.assign(content, { format: 'org.matrix.custom.html', formatted_body: html });
	return content;
}

/** The `m.relates_to` block for a reply, a thread message, or both. */
export function relation(target: Target, lastInThread?: string): Content | null {
	if (target.threadId) {
		const replyId = target.replyTo?.getId() ?? lastInThread ?? target.threadId;
		return {
			rel_type: THREAD_RELATION_TYPE.name,
			event_id: target.threadId,
			is_falling_back: !target.replyTo,
			'm.in_reply_to': { event_id: replyId }
		};
	}
	const replyId = target.replyTo?.getId();
	return replyId ? { 'm.in_reply_to': { event_id: replyId } } : null;
}

/** Sends a text message with markdown, mentions, and an optional reply or thread. */
export async function sendText(client: MatrixClient, target: Target, draft: Draft, lastInThread?: string) {
	const content = textContent(draft);
	const rel = relation(target, lastInThread);
	if (rel) content['m.relates_to'] = rel;
	await client.sendMessage(target.roomId, target.threadId ?? null, content as never);
}

/** Replaces a message's content with an `m.replace` edit. */
export async function editText(client: MatrixClient, event: MatrixEvent, draft: Draft) {
	const next = textContent(draft);
	const content: Content = {
		...next,
		body: `* ${next.body as string}`,
		'm.new_content': next,
		'm.relates_to': { rel_type: RelationType.Replace, event_id: event.getId() }
	};
	await client.sendMessage(event.getRoomId()!, event.threadRootId ?? null, content as never);
}

/** Deletes a message by redacting it. */
export async function deleteMessage(client: MatrixClient, event: MatrixEvent) {
	await client.redactEvent(event.getRoomId()!, event.getId()!);
}

/** Adds a reaction, or removes the user's own reaction with the same key. */
export async function toggleReaction(client: MatrixClient, event: MatrixEvent, key: string, mine?: MatrixEvent) {
	const roomId = event.getRoomId()!;
	if (mine?.getId()) {
		await client.redactEvent(roomId, mine.getId()!);
		return;
	}
	const content = { 'm.relates_to': { rel_type: RelationType.Annotation, event_id: event.getId(), key } };
	await client.sendEvent(roomId, event.threadRootId ?? null, EventType.Reaction, content as never);
}
