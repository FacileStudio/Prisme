import { escapeHtml } from './markdown';
import { handle } from './users';

export type Mentionable = { userId: string; name: string; agent: boolean };
export type Mention = { userId: string; name: string };
export type MentionRange = { start: number; end: number; mention: Mention };
export type Segment = { text: string; mention: Mention | null };

/** The mention of a user: their username, written `@pixel`, never their display name. */
export function mentionFor(userId: string): Mention {
	return { userId, name: handle(userId).replace(/^@/, '') };
}

const WORD =/[\p{L}\p{N}_]/u;

function matchesAt(text: string, at: number, name: string): boolean {
	if (!name || !text.startsWith(name, at + 1)) return false;
	const next = text[at + 1 + name.length];
	return next === undefined || !WORD.test(next);
}

/**
 * Finds every `@name` of a known mention in the text. An `@` glued to a word
 * (an email) never starts one, a name must end on a word boundary, and the
 * longest name wins; on a tie the earlier entry in `known` wins.
 */
export function findMentions(text: string, known: Mention[]): MentionRange[] {
	const byLength = [...known].sort((a, b) => b.name.length - a.name.length);
	const ranges: MentionRange[] = [];
	let at = text.indexOf('@');
	while (at !== -1) {
		const free = at === 0 || !WORD.test(text[at - 1]);
		const hit = free ? byLength.find((m) => matchesAt(text, at, m.name)) : undefined;
		const end = hit ? at + hit.name.length + 1 : at + 1;
		if (hit) ranges.push({ start: at, end, mention: hit });
		at = text.indexOf('@', end);
	}
	return ranges;
}

/** Splits text into plain runs and `@name` mention runs, for the composer highlight. */
export function segments(text: string, known: Mention[]): Segment[] {
	const out: Segment[] = [];
	let from = 0;
	for (const range of findMentions(text, known)) {
		if (range.start > from) out.push({ text: text.slice(from, range.start), mention: null });
		out.push({ text: text.slice(range.start, range.end), mention: range.mention });
		from = range.end;
	}
	if (from < text.length) out.push({ text: text.slice(from), mention: null });
	return out;
}

/** The mentions whose `@name` is in the text, in order of first appearance, one per user. */
export function presentMentions(text: string, mentions: Mention[]): Mention[] {
	const seen = new Set<string>();
	return findMentions(text, mentions)
		.map((range) => range.mention)
		.filter((m) => !seen.has(m.userId) && seen.add(m.userId));
}

/** Reports whether a message's `m.mentions` names the given user. */
export function mentionsUser(content: Record<string, unknown>, userId: string): boolean {
	const block = content['m.mentions'] as { user_ids?: unknown } | undefined;
	return Array.isArray(block?.user_ids) && block.user_ids.includes(userId);
}

/** The matrix.to permalink Matrix clients render as a user pill. */
export function permalink(userId: string): string {
	return `https://matrix.to/#/${userId}`;
}

/** The HTML pill for one mention. */
export function pill(m: Mention): string {
	return `<a href="${permalink(m.userId)}">${escapeHtml(`@${m.name}`)}</a>`;
}
