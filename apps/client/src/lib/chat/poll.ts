import { M_POLL_KIND_UNDISCLOSED, M_POLL_RESPONSE, M_POLL_START, type MatrixEvent } from 'matrix-js-sdk';

export const MAX_ANSWERS = 20;

export type PollAnswer = { id: string; text: string };
export type PollInfo = { question: string; answers: PollAnswer[]; maxSelections: number; disclosed: boolean };
export type PollVote = { sender: string; ts: number; answers: string[] };
export type PollTally = { counts: Map<string, number>; voters: number; mine: string[]; leader: number };

const START_KEYS = [M_POLL_START.name, M_POLL_START.altName] as string[];
const RESPONSE_KEYS = [M_POLL_RESPONSE.name, M_POLL_RESPONSE.altName] as string[];

/** Reports whether an event type is a poll start, stable or unstable. */
export function isPollStartType(type: string): boolean {
	return START_KEYS.includes(type);
}

function textOf(block: unknown): string {
	if (typeof block === 'string') return block;
	if (!block || typeof block !== 'object') return '';
	const record = block as Record<string, unknown>;
	for (const key of ['m.text', 'org.matrix.msc1767.text', 'm.message', 'org.matrix.msc1767.message', 'body']) {
		const value = record[key];
		if (typeof value === 'string') return value;
		if (Array.isArray(value)) {
			const plain = value.find((part) => !part?.mimetype || part.mimetype === 'text/plain') ?? value[0];
			if (typeof plain?.body === 'string') return plain.body;
		}
	}
	return '';
}

/** The poll a start event's content describes, or null when it is not a usable poll. */
export function pollOf(content: Record<string, unknown>): PollInfo | null {
	const key = START_KEYS.find((k) => content[k] && typeof content[k] === 'object');
	if (!key) return null;
	const start = content[key] as Record<string, unknown>;
	const question = textOf(start.question).trim();
	const raw = Array.isArray(start.answers) ? (start.answers as Record<string, unknown>[]) : [];
	const answers = raw
		.map((a) => ({ id: typeof a?.id === 'string' ? a.id : '', text: textOf(a).trim() }))
		.filter((a) => a.id && a.text)
		.slice(0, MAX_ANSWERS);
	if (!question || answers.length < 1) return null;
	const max = Number(start.max_selections);
	const kind = String(start.kind ?? '');
	return {
		question,
		answers,
		maxSelections: Number.isInteger(max) && max >= 1 ? Math.min(max, answers.length) : 1,
		disclosed: kind !== M_POLL_KIND_UNDISCLOSED.name && kind !== M_POLL_KIND_UNDISCLOSED.altName
	};
}

/** The answer ids a response event's content picks, before validation. */
export function responseAnswers(content: Record<string, unknown>): string[] {
	const key = RESPONSE_KEYS.find((k) => content[k] && typeof content[k] === 'object');
	if (!key) return [];
	const answers = (content[key] as { answers?: unknown }).answers;
	return Array.isArray(answers) ? answers.filter((a): a is string => typeof a === 'string') : [];
}

/**
 * Counts a poll the way MSC3381 says: each sender's latest response before the
 * end counts, only known answer ids count, and at most `maxSelections` of them.
 * A latest response with no valid answer spoils that sender's vote.
 */
export function tally(poll: PollInfo, votes: PollVote[], me: string, endTs?: number): PollTally {
	const known = new Set(poll.answers.map((a) => a.id));
	const latest = new Map<string, PollVote>();
	for (const vote of votes) {
		if (endTs !== undefined && vote.ts > endTs) continue;
		const prev = latest.get(vote.sender);
		if (!prev || vote.ts >= prev.ts) latest.set(vote.sender, vote);
	}
	const counts = new Map(poll.answers.map((a) => [a.id, 0]));
	let voters = 0;
	let mine: string[] = [];
	for (const [sender, vote] of latest) {
		const picks = [...new Set(vote.answers.filter((id) => known.has(id)))].slice(0, poll.maxSelections);
		if (sender === me) mine = picks;
		if (picks.length === 0) continue;
		voters += 1;
		for (const id of picks) counts.set(id, (counts.get(id) ?? 0) + 1);
	}
	return { counts, voters, mine, leader: Math.max(0, ...counts.values()) };
}

/** The poll a poll start event opens, if the event is one. */
export function pollInfo(event: MatrixEvent): PollInfo | null {
	return isPollStartType(event.getType()) ? pollOf(event.getContent()) : null;
}
