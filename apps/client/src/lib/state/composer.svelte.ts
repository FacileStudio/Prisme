import type { MatrixClient, Room } from 'matrix-js-sdk';
import { activeMention, filterMembers, insertMention } from '$lib/chat/complete';
import { mentionFor, type Mention, type Mentionable } from '$lib/chat/mentions';
import { isAgent, isBot } from '$lib/chat/users';
import { TypingNotifier } from '$lib/matrix/typing';
import { stripReplyFallback } from '$lib/chat/markdown';
import { editText, sendCapsule, sendText } from '$lib/matrix/messages';
import { sealText, type Expiry } from '$lib/chat/capsule';
import { composeMessage } from '$lib/chat/compose';
import { t } from '$lib/i18n/index.svelte';
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
	sealed = $state(false);
	expiry = $state<Expiry>('24h');
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

	/** Replaces the selection with text, such as an emoji, and returns the caret after it. */
	insert(text: string, start: number, end = start): number {
		this.text = this.text.slice(0, start) + text + this.text.slice(end);
		this.caret = start + text.length;
		this.typed();
		return this.caret;
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
		} else if (this.sealed) {
			await this.sendSealed(draft, message);
		} else {
			for (const file of this.files) await sendFile(this.client, this.room.roomId, file, this.threadId);
			const last = this.threadId ? this.room.getThread(this.threadId)?.lastReply()?.getId() : undefined;
			const target = { roomId: this.room.roomId, threadId: this.threadId, replyTo: draft.replyTo };
			if (message.text.trim()) await sendText(this.client, target, message, last);
		}
		this.reset();
		draft.clear();
	}

	private async sendSealed(draft: DraftState, message: { text: string; mentions: Mention[] }) {
		if (!message.text.trim()) return;
		const capsule = await sealText(message.text, this.expiry);
		const { userIds } = composeMessage(message.text, message.mentions);
		const last = this.threadId ? this.room.getThread(this.threadId)?.lastReply()?.getId() : undefined;
		const target = { roomId: this.room.roomId, threadId: this.threadId, replyTo: draft.replyTo };
		await sendCapsule(this.client, target, capsule, t('capsule.body'), userIds, last);
	}

	private reset() {
		this.text = '';
		this.caret = 0;
		this.files = [];
		this.mentions = [];
		this.sealed = false;
	}
}
