import type { MatrixClient, Room } from 'matrix-js-sdk';
import { activeMention, filterMembers, insertMention } from '$lib/chat/complete';
import { mentionFor, type Mention, type Mentionable } from '$lib/chat/mentions';
import { isAgent, isBot } from '$lib/chat/users';
import { TypingNotifier } from '$lib/matrix/typing';
import { stripReplyFallback } from '$lib/chat/markdown';
import { editText, sendText } from '$lib/matrix/messages';
import { sendFile } from '$lib/matrix/upload';
import type { DraftState } from './draft.svelte';

/** The members of a room as mention candidates, agents flagged, agents first. */
export function mentionables(room: Room, me: string, agents?: ReadonlyMap<string, unknown>): Mentionable[] {
	return room
		.getJoinedMembers()
		.filter((m) => m.userId !== me && !isBot(m.userId))
		.map((m) => ({ userId: m.userId, name: m.rawDisplayName || m.name, agent: isAgent(m.userId, agents) }));
}

/** One composer's text, mentions, completion and attachments. */
export class ComposerState {
	text = $state('');
	caret = $state(0);
	files = $state.raw<File[]>([]);
	selected = $state(0);
	sending = $state(false);
	dismissed = $state(-1);
	mentions = $state.raw<Mention[]>([]);
	readonly typing: TypingNotifier;

	constructor(
		private client: MatrixClient,
		private room: Room,
		private threadId: string | null
	) {
		this.typing = new TypingNotifier(client, room.roomId);
	}

	/** Everyone the text may name: the chosen mentions first, then the room members. */
	known(members: Mentionable[]): Mention[] {
		return [...this.mentions, ...members.map(({ userId }) => mentionFor(userId))];
	}

	/** The completion candidates for the `@query` under the caret. */
	completions(members: Mentionable[]): Mentionable[] {
		const at = activeMention(this.text, this.caret);
		return at && at.start !== this.dismissed ? filterMembers(members, at.query) : [];
	}

	/** Inserts a chosen member and returns the new caret position. */
	choose(member: Mentionable): number {
		const at = activeMention(this.text, this.caret);
		if (!at) return this.caret;
		const mention = mentionFor(member.userId);
		const next = insertMention(this.text, at, this.caret, mention.name);
		this.mentions = [...this.mentions, mention];
		this.text = next.text;
		this.caret = next.caret;
		this.selected = 0;
		return next.caret;
	}

	/** Reports a keystroke to the typing notifier. */
	typed() {
		this.typing.keystroke(this.text.trim().length > 0);
	}

	/** Loads a message into the composer for editing. */
	loadEdit(draft: DraftState) {
		const event = draft.editing;
		if (!event) return;
		const content = event.getContent();
		this.text = stripReplyFallback(String(content.body ?? ''));
		const ids = (content['m.mentions'] as { user_ids?: string[] } | undefined)?.user_ids ?? [];
		this.mentions = ids.map(mentionFor);
	}

	/** Sends the attachments and the text, or saves an edit. */
	async submit(draft: DraftState, members: Mentionable[]) {
		this.typing.stop();
		const message = { text: this.text, mentions: this.known(members) };
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
	}
}
