import type { Locale } from '@facile/muse';

/** One emoji of the picker, with the variants of its skin tones in order, light to dark. */
export type Emoji = { unicode: string; label: string; tags: string[]; group: number; order: number; skins: string[] };

/** The emojibase groups the picker shows, in order. Group 2 holds bare skin-tone components. */
export const EMOJI_GROUPS = [0, 1, 3, 4, 5, 6, 7, 8, 9] as const;
export type EmojiGroup = (typeof EMOJI_GROUPS)[number];

/** The five skin-tone modifiers, light to dark. Tone 0 is the default yellow. */
export const SKIN_TONES = ['1F3FB', '1F3FC', '1F3FD', '1F3FE', '1F3FF'] as const;

type Raw = { unicode: string; label: string; hexcode: string; tags?: string[]; group?: number; order?: number; skins?: Raw[] };

const loaders: Record<Locale, () => Promise<{ default: unknown }>> = {
	en: () => import('emojibase-data/en/compact.json'),
	fr: () => import('emojibase-data/fr/compact.json'),
	es: () => import('emojibase-data/es/compact.json'),
	de: () => import('emojibase-data/de/compact.json')
};

const cache = new Map<Locale, Promise<Emoji[]>>();

/** The picker's emojis in the given locale, loaded once and only when the picker first opens. */
export function loadEmoji(locale: Locale): Promise<Emoji[]> {
	let pending = cache.get(locale);
	if (!pending) {
		pending = loaders[locale]().then((module) => toEmoji(module.default as Raw[]));
		pending.catch(() => cache.delete(locale));
		cache.set(locale, pending);
	}
	return pending;
}

/** Keeps the emojis a picker can show and sorts them the way Unicode orders them. */
export function toEmoji(raw: Raw[]): Emoji[] {
	return raw
		.filter((e) => e.group !== undefined && e.group !== 2 && e.order !== undefined)
		.map((e) => ({
			unicode: e.unicode,
			label: e.label,
			tags: e.tags ?? [],
			group: e.group!,
			order: e.order!,
			skins: SKIN_TONES.map((tone) => e.skins?.find((s) => singleTone(s.hexcode) === tone)?.unicode ?? '')
		}))
		.sort((a, b) => a.order - b.order);
}

function singleTone(hexcode: string): string | null {
	const tones = hexcode.split('-').filter((part) => (SKIN_TONES as readonly string[]).includes(part));
	return tones.length === 1 || new Set(tones).size === 1 ? tones[0] : null;
}

/** The emoji in the chosen tone: 0 is the default, 1 to 5 go light to dark. */
export function withTone(emoji: Emoji, tone: number): string {
	return (tone > 0 && emoji.skins[tone - 1]) || emoji.unicode;
}

function fold(text: string): string {
	return text.normalize('NFD').replace(/\p{M}/gu, '').toLowerCase();
}

/** Emojis whose name or tags start with every word of the query, best matches first. */
export function searchEmoji(list: Emoji[], query: string, limit = 96): Emoji[] {
	const words = fold(query).split(/\s+/).filter(Boolean);
	if (words.length === 0) return [];
	const scored: { emoji: Emoji; score: number }[] = [];
	for (const emoji of list) {
		const label = fold(emoji.label);
		const terms = [...label.split(/[\s:,-]+/), ...emoji.tags.map(fold)];
		if (!words.every((w) => terms.some((term) => term.startsWith(w)))) continue;
		const score = label === words.join(' ') ? 0 : label.startsWith(words[0]) ? 1 : 2;
		scored.push({ emoji, score });
	}
	return scored
		.sort((a, b) => a.score - b.score || a.emoji.order - b.emoji.order)
		.slice(0, limit)
		.map((entry) => entry.emoji);
}
