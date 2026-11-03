import { RoomMemberEvent, type MatrixClient, type MatrixEvent, type Room, type RoomMember } from 'matrix-js-sdk';

const IDLE = 4000;
const REFRESH = 3000;
const SERVER_TIMEOUT = 6000;

/**
 * Sends typing notices for one room: on while keys arrive (refreshed before the
 * server timeout), off after four idle seconds, and off on demand (send, blur,
 * leaving the room).
 */
export class TypingNotifier {
	private sentAt = 0;
	private idle: ReturnType<typeof setTimeout> | undefined;

	constructor(
		private client: MatrixClient,
		private roomId: string
	) {}

	/** Records a keystroke; `writing` is false once the draft is empty. */
	keystroke(writing: boolean) {
		if (!writing) return this.stop();
		clearTimeout(this.idle);
		this.idle = setTimeout(() => this.stop(), IDLE);
		const now = Date.now();
		if (now - this.sentAt < REFRESH) return;
		this.sentAt = now;
		this.send(true);
	}

	/** Clears the notice if one is on. */
	stop() {
		clearTimeout(this.idle);
		this.idle = undefined;
		if (!this.sentAt) return;
		this.sentAt = 0;
		this.send(false);
	}

	private send(typing: boolean) {
		void this.client.sendTyping(this.roomId, typing, typing ? SERVER_TIMEOUT : 0).catch(() => undefined);
	}
}

/** The members of a room typing right now, the current user left out. */
export function typingMembers(room: Room, me: string): RoomMember[] {
	return room.getMembers().filter((m) => m.typing && m.userId !== me);
}

/** Calls `onChange` with the typing members of a room each time a `RoomMember.typing` event changes them. */
export function watchTyping(client: MatrixClient, room: Room, onChange: (members: RoomMember[]) => void): () => void {
	const me = client.getUserId() ?? '';
	const handler = (_event: MatrixEvent, member: RoomMember) => {
		if (member.roomId === room.roomId) onChange(typingMembers(room, me));
	};
	client.on(RoomMemberEvent.Typing, handler);
	onChange(typingMembers(room, me));
	return () => client.off(RoomMemberEvent.Typing, handler);
}
