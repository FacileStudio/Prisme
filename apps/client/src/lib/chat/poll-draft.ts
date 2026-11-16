import { M_POLL_KIND_DISCLOSED, M_POLL_KIND_UNDISCLOSED, M_POLL_START } from 'matrix-js-sdk';
import { MAX_ANSWERS } from './poll';

export const MIN_ANSWERS = 2;
export const MAX_QUESTION = 340;
export const MAX_ANSWER = 340;

export type PollDraft = { question: string; answers: string[]; multiple: boolean; disclosed: boolean };

const M_TEXT_KEY = 'm.text';

/** The next selection after a click: single choice replaces, multiple choice toggles up to the limit. */
export function toggleAnswer(current: string[], id: string, maxSelections: number): string[] {
	if (maxSelections <= 1) return current.length === 1 && current[0] === id ? [] : [id];
	if (current.includes(id)) return current.filter((x) => x !== id);
	return current.length >= maxSelections ? current : [...current, id];
}

/** Whether a draft can be sent: a question and at least two distinct non-empty answers. */
export function draftReady(draft: PollDraft): boolean {
	const answers = draft.answers.map((a) => a.trim()).filter(Boolean);
	return !!draft.question.trim() && answers.length >= MIN_ANSWERS && new Set(answers).size === answers.length;
}

/** The poll start content for a draft, with the plain-text fallback other clients show. */
export function pollContent(draft: PollDraft, newId: () => string): Record<string, unknown> {
	const question = draft.question.trim().slice(0, MAX_QUESTION);
	const answers = draft.answers
		.map((a) => a.trim().slice(0, MAX_ANSWER))
		.filter(Boolean)
		.slice(0, MAX_ANSWERS)
		.map((text) => ({ id: newId(), [M_TEXT_KEY]: text }));
	const fallback = `${question}\n${answers.map((a, i) => `${i + 1}. ${a[M_TEXT_KEY]}`).join('\n')}`;
	return {
		[M_POLL_START.name]: {
			question: { [M_TEXT_KEY]: question },
			kind: draft.disclosed ? M_POLL_KIND_DISCLOSED.name : M_POLL_KIND_UNDISCLOSED.name,
			max_selections: draft.multiple ? answers.length : 1,
			answers
		},
		[M_TEXT_KEY]: fallback,
		body: fallback
	};
}


/** The share of voters that picked an answer, as a whole percentage. */
export function percent(count: number, voters: number): number {
	return voters === 0 ? 0 : Math.round((count / voters) * 100);
}
