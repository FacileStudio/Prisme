import { describe, expect, it } from 'vitest';
import { searchEmoji, toEmoji, withTone } from './emoji';

const raw = [
	{ unicode: '🇦', label: 'regional indicator A', hexcode: '1F1E6' },
	{ unicode: '🏻', label: 'light skin tone', hexcode: '1F3FB', group: 2, order: 1 },
	{
		unicode: '👋',
		label: 'main qui salue',
		hexcode: '1F44B',
		tags: ['bonjour', 'salut'],
		group: 1,
		order: 189,
		skins: [
			{ unicode: '👋🏻', label: 'x', hexcode: '1F44B-1F3FB' },
			{ unicode: '👋🏿', label: 'x', hexcode: '1F44B-1F3FF' }
		]
	},
	{ unicode: '😀', label: 'visage très souriant', hexcode: '1F600', tags: ['sourire'], group: 0, order: 1 },
	{ unicode: '🐱', label: 'tête de chat', hexcode: '1F431', tags: ['chat'], group: 3, order: 600 }
];

describe('emoji', () => {
	const list = toEmoji(raw);

	it('drops components and ungrouped entries and keeps Unicode order', () => {
		expect(list.map((e) => e.unicode)).toEqual(['😀', '👋', '🐱']);
	});

	it('maps single-tone variants to their slot', () => {
		const wave = list[1];
		expect(withTone(wave, 0)).toBe('👋');
		expect(withTone(wave, 1)).toBe('👋🏻');
		expect(withTone(wave, 3)).toBe('👋');
		expect(withTone(wave, 5)).toBe('👋🏿');
	});

	it('matches word prefixes of labels and tags, ignoring accents', () => {
		expect(searchEmoji(list, 'tres').map((e) => e.unicode)).toEqual(['😀']);
		expect(searchEmoji(list, 'chat').map((e) => e.unicode)).toEqual(['🐱']);
		expect(searchEmoji(list, 'sal').map((e) => e.unicode)).toEqual(['👋']);
		expect(searchEmoji(list, 'main sal').map((e) => e.unicode)).toEqual(['👋']);
		expect(searchEmoji(list, 'zzz')).toEqual([]);
		expect(searchEmoji(list, '  ')).toEqual([]);
	});
});
