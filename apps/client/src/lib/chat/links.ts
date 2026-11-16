import { prisme, type LinkPreview } from '$lib/api/prisme';

export type SuiteApp = 'feuilles' | 'docs' | 'echo' | 'clap' | 'diapo' | 'nuage';
export type SuiteLink = { url: string; app: SuiteApp };

const HOSTS: Record<string, SuiteApp> = {
	'feuilles.facile.studio': 'feuilles',
	'docs.facile.studio': 'docs',
	'echo.facile.studio': 'echo',
	'clap.facile.studio': 'clap',
	'diapo.facile.studio': 'diapo',
	'nuage.facile.studio': 'nuage'
};

const URL_PATTERN = /https:\/\/[^\s<>"'`]+/gi;
const TRAILING = /[.,;:!?)\]}]+$/;

/**
 * The links to other Facile apps in a message, in order and without repeats.
 * Capsule links are never listed: previewing one would burn it.
 */
export function suiteLinks(text: string, limit = 3): SuiteLink[] {
	const found: SuiteLink[] = [];
	const seen = new Set<string>();
	for (const match of text.matchAll(URL_PATTERN)) {
		const link = suiteLink(match[0].replace(TRAILING, ''));
		if (!link || seen.has(link.url)) continue;
		seen.add(link.url);
		found.push(link);
		if (found.length === limit) break;
	}
	return found;
}

function suiteLink(raw: string): SuiteLink | null {
	try {
		const url = new URL(raw);
		const app = HOSTS[url.hostname];
		if (url.protocol !== 'https:' || url.port || url.username || !app) return null;
		return { url: url.href, app };
	} catch {
		return null;
	}
}

/** A link without its fragment, which may carry a key the server must never see. */
export function withoutFragment(raw: string): string {
	const url = new URL(raw);
	url.hash = '';
	return url.href;
}

const previews = new Map<string, Promise<LinkPreview | null>>();

/** The card data for a link, fetched once per page load. */
export function linkPreview(raw: string): Promise<LinkPreview | null> {
	const key = withoutFragment(raw);
	let pending = previews.get(key);
	if (!pending) {
		pending = prisme.unfurl(key).catch(() => null);
		previews.set(key, pending);
	}
	return pending;
}
