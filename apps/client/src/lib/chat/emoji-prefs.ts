const RECENT_KEY = 'prisme.emoji.recent';
const TONE_KEY = 'prisme.emoji.tone';
const RECENT_MAX = 24;

/** The emojis this browser picked last, newest first. */
export function recentEmoji(): string[] {
	try {
		const value = JSON.parse(localStorage.getItem(RECENT_KEY) ?? '[]');
		return Array.isArray(value) ? value.filter((v): v is string => typeof v === 'string').slice(0, RECENT_MAX) : [];
	} catch {
		return [];
	}
}

/** Puts an emoji at the front of the recent list. */
export function rememberEmoji(unicode: string): string[] {
	const next = [unicode, ...recentEmoji().filter((e) => e !== unicode)].slice(0, RECENT_MAX);
	try {
		localStorage.setItem(RECENT_KEY, JSON.stringify(next));
	} catch {
		return next;
	}
	return next;
}

/** The skin tone this browser chose last, 0 to 5. */
export function savedTone(): number {
	try {
		const tone = Number(localStorage.getItem(TONE_KEY));
		return Number.isInteger(tone) && tone >= 0 && tone <= 5 ? tone : 0;
	} catch {
		return 0;
	}
}

/** Remembers the chosen skin tone. */
export function saveTone(tone: number): void {
	try {
		localStorage.setItem(TONE_KEY, String(tone));
	} catch {
		return;
	}
}
