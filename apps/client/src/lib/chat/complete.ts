import type { Mentionable } from './mentions';

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
