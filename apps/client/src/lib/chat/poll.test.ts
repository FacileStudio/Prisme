import { describe, expect, it } from 'vitest';
import { pollOf, responseAnswers, tally, type PollInfo } from './poll';
import { draftReady, percent, pollContent, toggleAnswer } from './poll-draft';

const poll: PollInfo = {
	question: 'Lunch?',
	answers: [
		{ id: 'a', text: 'Pizza' },
		{ id: 'b', text: 'Sushi' },
		{ id: 'c', text: 'Salad' }
	],
	maxSelections: 1,
	disclosed: true
};

describe('poll tally', () => {
	it('counts only each sender latest valid vote', () => {
		const result = tally(
			poll,
			[
				{ sender: '@ann', ts: 1, answers: ['a'] },
				{ sender: '@ann', ts: 5, answers: ['b'] },
				{ sender: '@bob', ts: 2, answers: ['b'] },
				{ sender: '@cy', ts: 3, answers: ['zzz'] }
			],
			'@ann'
		);
		expect(result.counts.get('a')).toBe(0);
		expect(result.counts.get('b')).toBe(2);
		expect(result.voters).toBe(2);
		expect(result.mine).toEqual(['b']);
		expect(result.leader).toBe(2);
	});

	it('ignores votes after the end and caps selections', () => {
		const multi = { ...poll, maxSelections: 2 };
		const result = tally(
			multi,
			[
				{ sender: '@ann', ts: 1, answers: ['a', 'b', 'c'] },
				{ sender: '@bob', ts: 9, answers: ['c'] }
			],
			'@me',
			5
		);
		expect([...result.counts.values()]).toEqual([1, 1, 0]);
		expect(result.voters).toBe(1);
	});

	it('lets an empty latest response withdraw a vote', () => {
		const result = tally(poll, [
			{ sender: '@ann', ts: 1, answers: ['a'] },
			{ sender: '@ann', ts: 2, answers: [] }
		], '@ann');
		expect(result.voters).toBe(0);
		expect(result.mine).toEqual([]);
	});
});

describe('poll content', () => {
	it('round trips a draft through the start content', () => {
		let n = 0;
		const content = pollContent({ question: ' Lunch? ', answers: ['Pizza', '', 'Sushi'], multiple: false, disclosed: false }, () => `id${n++}`);
		const parsed = pollOf(content);
		expect(parsed?.question).toBe('Lunch?');
		expect(parsed?.answers).toEqual([
			{ id: 'id0', text: 'Pizza' },
			{ id: 'id1', text: 'Sushi' }
		]);
		expect(parsed?.maxSelections).toBe(1);
		expect(parsed?.disclosed).toBe(false);
		expect(String(content.body)).toContain('1. Pizza');
	});

	it('reads unstable polls with text arrays', () => {
		const parsed = pollOf({
			'org.matrix.msc3381.poll.start': {
				question: { 'org.matrix.msc1767.text': 'Q' },
				kind: 'org.matrix.msc3381.poll.disclosed',
				max_selections: 9,
				answers: [{ id: 'x', 'm.text': [{ body: 'Yes' }] }, { id: 'y', 'org.matrix.msc1767.text': 'No' }]
			}
		});
		expect(parsed?.answers.map((a) => a.text)).toEqual(['Yes', 'No']);
		expect(parsed?.maxSelections).toBe(2);
		expect(parsed?.disclosed).toBe(true);
	});

	it('rejects content that is not a poll', () => {
		expect(pollOf({ body: 'hi' })).toBeNull();
		expect(responseAnswers({ 'm.poll.response': { answers: ['a', 3] } })).toEqual(['a']);
	});

	it('needs two distinct answers and a question', () => {
		expect(draftReady({ question: 'Q', answers: ['a', 'a'], multiple: false, disclosed: true })).toBe(false);
		expect(draftReady({ question: '', answers: ['a', 'b'], multiple: false, disclosed: true })).toBe(false);
		expect(draftReady({ question: 'Q', answers: ['a', ' b ', ''], multiple: false, disclosed: true })).toBe(true);
	});

	it('toggles single and multiple choice', () => {
		expect(toggleAnswer(['a'], 'b', 1)).toEqual(['b']);
		expect(toggleAnswer(['a'], 'a', 1)).toEqual([]);
		expect(toggleAnswer(['a', 'b'], 'c', 2)).toEqual(['a', 'b']);
		expect(toggleAnswer(['a', 'b'], 'a', 2)).toEqual(['b']);
		expect(percent(1, 3)).toBe(33);
		expect(percent(0, 0)).toBe(0);
	});
});
