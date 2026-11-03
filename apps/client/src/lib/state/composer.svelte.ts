import type { MatrixClient, Room } from 'matrix-js-sdk';
import { activeMention, filterMembers, insertMention, type Mention, type Mentionable } from '$lib/chat/mentions';
import { isAgent, isBot } from '$lib/chat/users';
import { stripReplyFallback } from '$lib/chat/markdown';
import { editText, sendText } from '$lib/matrix/messages';
import { sendFile } from '$lib/matrix/upload';
import type { DraftState } from './draft.svelte';

/** The members of a room as mention candidates, agents flagged. */
export function mentionables(room: Room, me: string): Mentionable[] {
	return room
		.getJoinedMembers()
		.filter((m) => m.userId !== me && !isBot(m.userId))
		.map((m) => ({ userId: m.userId, name: m.name, agent: isAgent(m.userId) }));
}

/** One composer's text, mentions, completion and attachments. */
export class ComposerState {
	text = $state('');
	caret = $state(0);
	files = $state.raw<File[]>([]);
	selected = $state(0);
	sending = $state(false);
	mentions: Mention[] = [];
	private typingAt = 0;

	constructor(
		private client: MatrixClient,
		private room: Room,
		private threadId: string | null
	) {}

	/** The completion candidates for the `@query` under the caret. */
	completions(members: Mentionable[]): Mentionable[] {
		const at = activeMention(this.text, this.caret);
		return at ? filterMembers(members, at.query) : [];
	}

	/** Inserts a chosen member and returns the new caret position. */
	choose(member: Mentionable): number {
		const at = activeMention(this.text, this.caret);
		if (!at) return this.caret;
		const next = insertMention(this.text, at, this.caret, member.name);
		this.mentions = [...this.mentions, { userId: member.userId, name: member.name }];
		this.text = next.text;
		this.caret = next.caret;
		this.selected = 0;
		return next.caret;
	}

	/** Sends typing notifications at most every four seconds while writing. */
	typed() {
		const now = Date.now();
		const typing = this.text.trim().length > 0;
		if (typing && now - this.typingAt < 4000) return;
		this.typingAt = typing ? now : 0;
		void this.client.sendTyping(this.room.roomId, typing, 6000).catch(() => undefined);
	}

	/** Loads a message into the composer for editing. */
	loadEdit(draft: DraftState) {
		const event = draft.editing;
		if (!event) return;
		const content = event.getContent();
		this.text = stripReplyFallback(String(content.body ?? ''));
		const ids = (content['m.mentions'] as { user_ids?: string[] } | undefined)?.user_ids ?? [];
		this.mentions = ids.map((id) => ({ userId: id, name: this.room.getMember(id)?.name ?? id }));
	}

	/** Sends the attachments and the text, or saves an edit. */
	async submit(draft: DraftState) {
		const message = { text: this.text, mentions: this.mentions };
		if (draft.editing) {
			if (message.text.trim()) await editText(this.client, draft.editing, message);
		} else {
			for (const file of this.files) await sendFile(this.client, this.room.roomId, file, this.threadId);
			const last = this.threadId ? this.room.getThread(this.threadId)?.lastReply()?.getId() : undefined;
			const target = { roomId: this.room.roomId, threadId: this.threadId, replyTo: draft.replyTo };
			if (message.text.trim()) await sendText(this.client, target, message, last);
		}
		this.reset();
		draft.clear();
	}

	private reset() {
		this.text = '';
		this.caret = 0;
		this.files = [];
		this.mentions = [];
		this.typingAt = 0;
		void this.client.sendTyping(this.room.roomId, false, 0).catch(() => undefined);
	}
}
