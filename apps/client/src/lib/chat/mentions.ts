import { escapeHtml } from './markdown';

export type Mentionable = { userId: string; name: string; agent: boolean };
export type Mention = { userId: string; name: string };
export type MentionQuery = { start: number; query: string };

/** Finds the `@query` the caret sits in, if any. */
export function activeMention(text: string, caret: number): MentionQuery | null {
	const before = text.slice(0, caret);
	const match = /(^|\s)@([^\s@]*)$/.exec(before);
	if (!match) return null;
	return { start: caret - match[2].length - 1, query: match[2] };
}

function score(member: Mentionable, query: string): number {
	const name = member.name.toLowerCase();
	const local = member.userId.slice(1).toLowerCase();
	if (name.startsWith(query) || local.startsWith(query)) return 0;
	if (name.includes(query) || local.includes(query)) return 1;
	return -1;
}

/** Ranks members for the completion menu: prefix matches first, then substrings. */
export function filterMembers(members: Mentionable[], query: string, limit = 8): Mentionable[] {
	const q = query.toLowerCase();
	return members
		.map((member) => ({ member, rank: score(member, q) }))
		.filter((entry) => entry.rank >= 0)
		.sort((a, b) => a.rank - b.rank || a.member.name.localeCompare(b.member.name))
		.slice(0, limit)
		.map((entry) => entry.member);
}

/** Replaces the active `@query` with the chosen display name. */
export function insertMention(text: string, at: MentionQuery, caret: number, name: string) {
	const insert = `@${name} `;
	const next = text.slice(0, at.start) + insert + text.slice(caret);
	return { text: next, caret: at.start + insert.length };
}

/** Keeps only the mentions whose `@name` is still in the text, without duplicates. */
export function presentMentions(text: string, mentions: Mention[]): Mention[] {
	const seen = new Set<string>();
	return mentions.filter((m) => {
		if (seen.has(m.userId) || !text.includes(`@${m.name}`)) return false;
		seen.add(m.userId);
		return true;
	});
}

/** The matrix.to permalink Matrix clients render as a user pill. */
export function permalink(userId: string): string {
	return `https://matrix.to/#/${userId}`;
}

/** The HTML pill for one mention. */
export function pill(m: Mention): string {
	return `<a href="${permalink(m.userId)}">${escapeHtml(`@${m.name}`)}</a>`;
}
