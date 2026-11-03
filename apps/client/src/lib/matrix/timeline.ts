import { EventType, RelationType, type MatrixEvent, type Room } from 'matrix-js-sdk';
import { continues } from '$lib/chat/unread';

export type Item = { event: MatrixEvent; grouped: boolean; day: string | null };
export type Reaction = { key: string; count: number; mine?: MatrixEvent; senders: string[] };
export type CallInfo = { call_id: string; join_url: string; started_by: string };

/** Reports whether an event is a message the timeline shows, deleted ones included. */
export function isMessage(event: MatrixEvent): boolean {
	if (event.getType() !== EventType.RoomMessage) return false;
	return event.getRelation()?.rel_type !== RelationType.Replace;
}

function dayLabel(ts: number): string {
	return new Date(ts).toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: 'numeric' });
}

/** Turns raw events into timeline rows with sender grouping and day breaks. */
export function timelineItems(events: MatrixEvent[]): Item[] {
	const items: Item[] = [];
	let prev: MatrixEvent | undefined;
	const unique = [...new Map(events.filter(isMessage).map((e) => [e.getId() ?? e.getTxnId(), e])).values()];
	for (const event of unique) {
		const day = dayLabel(event.getTs());
		const newDay = !prev || dayLabel(prev.getTs()) !== day;
		const here = { sender: event.getSender() ?? '', ts: event.getTs() };
		const before = prev && !prev.isRedacted() ? { sender: prev.getSender() ?? '', ts: prev.getTs() } : undefined;
		items.push({ event, grouped: !newDay && continues(before, here), day: newDay ? day : null });
		prev = event;
	}
	return items;
}

/** The reactions on an event, grouped by key, with the user's own reaction. */
export function reactionsFor(room: Room, event: MatrixEvent, me: string): Reaction[] {
	const relations = room.relations.getChildEventsForEvent(event.getId()!, RelationType.Annotation, EventType.Reaction);
	const sorted = relations?.getSortedAnnotationsByKey() ?? [];
	return sorted.map(([key, events]) => {
		const list = [...events].filter((e) => !e.isRedacted());
		return {
			key,
			count: list.length,
			mine: list.find((e) => e.getSender() === me),
			senders: list.map((e) => e.getSender() ?? '')
		};
	}).filter((r) => r.count > 0);
}

/** The id of the message this one replies to, ignoring thread fallbacks. */
export function replyTarget(event: MatrixEvent): string | null {
	const rel = event.getWireContent()['m.relates_to'] as
		| { rel_type?: string; is_falling_back?: boolean; 'm.in_reply_to'?: { event_id?: string } }
		| undefined;
	if (!rel?.['m.in_reply_to']?.event_id) return null;
	if (rel.rel_type === 'm.thread' && rel.is_falling_back !== false) return null;
	return rel['m.in_reply_to'].event_id;
}

/** The Prisme call block of a call announcement, if the message is one. */
export function callInfo(event: MatrixEvent): CallInfo | null {
	const call = event.getContent()['studio.facile.call'] as CallInfo | undefined;
	return call?.call_id ? call : null;
}

/** Users other than the sender and the viewer whose read receipt sits on this event. */
export function seenBy(room: Room, event: MatrixEvent, me: string): string[] {
	const sender = event.getSender();
	return [...new Set(room.getUsersReadUpTo(event))].filter((id) => id !== me && id !== sender);
}
