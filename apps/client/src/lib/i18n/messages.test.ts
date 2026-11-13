import { describe, expect, it } from 'vitest';
import { catalogs } from './messages';

const locales = ['fr', 'es', 'de'] as const;

function placeholders(message: unknown): string {
	const texts = typeof message === 'string' ? [message] : Object.values(message as Record<string, string>);
	const names = texts.flatMap((text) => [...text.matchAll(/\{(\w+)\}/g)].map((match) => match[1]));
	return [...new Set(names)].sort().join(',');
}

describe('app messages', () => {
	it('carry every key in the four locales with the same placeholders', () => {
		const problems: string[] = [];
		for (const [key, english] of Object.entries(catalogs.en)) {
			for (const locale of locales) {
				const message = (catalogs[locale] as Record<string, unknown>)[key];
				if (!message) problems.push(`${locale} ${key}: empty`);
				else if (placeholders(message) !== placeholders(english)) problems.push(`${locale} ${key}: placeholders differ`);
			}
		}
		expect(problems).toEqual([]);
	});

	it('translates and falls back to English for an unknown locale', async () => {
		const { createTranslator } = await import('@facile/muse');
		const t = createTranslator(catalogs, 'fr');
		expect(t('nav.messages')).toBe('Messages');
		expect(t('space.channelCount', { count: 3 })).toBe('3 canaux');
		expect(t('typing.two', { a: 'Ann', b: 'Bob' })).toBe('Ann et Bob écrivent…');
	}, 30_000);
});
