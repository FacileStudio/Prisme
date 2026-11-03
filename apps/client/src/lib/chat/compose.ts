import { escapeHtml, renderMarkdown } from './markdown';
import { findMentions, pill, presentMentions, type Mention } from './mentions';

export type Composed = { body: string; html: string | null; userIds: string[] };

const token = (index: number) => `prismemention${index}x`;

function tokenize(text: string, known: Mention[], present: Mention[]): string {
	const ranges = findMentions(text, known);
	return ranges.reduceRight((out, range) => {
		const index = present.findIndex((m) => m.userId === range.mention.userId);
		return out.slice(0, range.start) + token(index) + out.slice(range.end);
	}, text);
}

/**
 * Turns composer text into a Matrix message: the plain body, the HTML body when
 * markdown or mentions need one, and the user ids for `m.mentions`. `known` lists
 * the people the text may name, the explicitly chosen ones first.
 */
export function composeMessage(text: string, known: Mention[]): Composed {
	const body = text.trim();
	const present = presentMentions(body, known);
	const tokenized = tokenize(body, known, present);
	let html = renderMarkdown(tokenized);
	if (html === null && present.length > 0) html = escapeHtml(tokenized).replace(/\n/g, '<br>');
	if (html !== null) {
		html = present.reduce((out, m, index) => out.split(token(index)).join(pill(m)), html);
	}
	return { body, html, userIds: present.map((m) => m.userId) };
}
