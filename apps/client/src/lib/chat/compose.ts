import { escapeHtml, renderMarkdown } from './markdown';
import { pill, presentMentions, type Mention } from './mentions';

export type Composed = { body: string; html: string | null; userIds: string[] };

const token = (index: number) => `prismemention${index}x`;

function tokenize(text: string, mentions: Mention[]): string {
	return mentions
		.map((m, index) => ({ m, index }))
		.sort((a, b) => b.m.name.length - a.m.name.length)
		.reduce((out, { m, index }) => out.split(`@${m.name}`).join(token(index)), text);
}

/**
 * Turns composer text into a Matrix message: the plain body, the HTML body when
 * markdown or mentions need one, and the user ids for `m.mentions`.
 */
export function composeMessage(text: string, mentions: Mention[]): Composed {
	const body = text.trim();
	const present = presentMentions(body, mentions);
	const tokenized = tokenize(body, present);
	let html = renderMarkdown(tokenized);
	if (html === null && present.length > 0) html = escapeHtml(tokenized).replace(/\n/g, '<br>');
	if (html !== null) {
		html = present.reduce((out, m, index) => out.split(token(index)).join(pill(m)), html);
	}
	return { body, html, userIds: present.map((m) => m.userId) };
}
