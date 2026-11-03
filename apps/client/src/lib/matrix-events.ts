import type { ChatMessage, RoomMember } from './types';

export type RawMatrixEvent = {
	event_id: string;
	type: string;
	sender: string;
	origin_server_ts: number;
	content: Record<string, any>;
	unsigned?: Record<string, any>;
};

export type JoinedRoomData = {
	timeline?: { events?: RawMatrixEvent[] };
	state?: { events?: RawMatrixEvent[] };
	ephemeral?: { events?: RawMatrixEvent[] };
};

export type JoinedRoomsData = Record<string, JoinedRoomData>;

export type SyncResponse = {
	next_batch: string;
	rooms?: {
		join?: JoinedRoomsData;
	};
};

export function parseChatMessage(roomId: string, event: RawMatrixEvent, isAgent: boolean): ChatMessage {
	const relatesTo = event.content?.['m.relates_to'];
	const isThread = relatesTo?.rel_type === 'm.thread';
	const threadRootId = isThread ? relatesTo?.event_id : undefined;
	const inReplyTo = relatesTo?.['m.in_reply_to']?.event_id;

	return {
		id: event.event_id,
		roomId,
		sender: event.sender,
		senderName: event.sender.split(':')[0].slice(1),
		isAgent,
		body: event.content?.body || '',
		timestamp: event.origin_server_ts,
		call: event.content?.['studio.facile.call'],
		replyToId: inReplyTo,
		threadRootId,
		reactions: {}
	};
}

export function parseMember(event: RawMatrixEvent, isAgent: boolean): RoomMember | null {
	if (event.type !== 'm.room.member') return null;
	const userId = event.content?.membership === 'join' ? event.sender : null;
	if (!userId) return null;
	return {
		userId,
		displayName: event.content?.displayname || userId.split(':')[0].slice(1),
		avatarUrl: event.content?.avatar_url,
		isAgent
	};
}
