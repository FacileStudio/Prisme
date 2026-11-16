import { glyphs } from '$lib/glyphs';

export const NUAGE_KEY = 'studio.facile.nuage';

/** The studio.facile.nuage block of a message sharing a Nuage file. The file stays on Nuage. */
export type NuageInfo = { file_id: number; name: string; size: number; mime: string; url: string };

/** What Prisme's Nuage share route answers for one file. */
export type NuageShared = { file_id: number; name: string; mime_type: string; size: number; url: string };

const SHARE_PATH = /^\/s\/[A-Za-z0-9_-]+$/;

/** The message block for a shared file, or null unless its link is a Nuage share page. */
export function nuageBlock(shared: NuageShared, origin = 'https://nuage.facile.studio'): NuageInfo | null {
	try {
		const url = new URL(shared.url);
		if (url.origin !== origin || !SHARE_PATH.test(url.pathname) || url.search || url.hash) return null;
	} catch {
		return null;
	}
	return { file_id: shared.file_id, name: shared.name, size: shared.size, mime: shared.mime_type, url: shared.url };
}

/** The Nuage block of raw message content, when it names a file and its link. */
export function nuageOf(content: Record<string, unknown>): NuageInfo | null {
	const block = content[NUAGE_KEY] as NuageInfo | undefined;
	if (typeof block?.name !== 'string' || typeof block.url !== 'string') return null;
	return block.url.startsWith('https://') ? block : null;
}

/** The plain body other Matrix clients show: the file name, then its link. */
export function nuageBody(info: NuageInfo): string {
	return `${info.name} ${info.url}`;
}

const KINDS: [RegExp, string][] = [
	[/^image\//, glyphs.camera],
	[/^video\//, glyphs.play],
	[/^audio\//, glyphs.volume]
];

/** The icon for a file's mime type, falling back to a plain file. */
export function fileGlyph(mime: string): string {
	return KINDS.find(([pattern]) => pattern.test(mime))?.[1] ?? glyphs.file;
}

const UNITS = ['B', 'KB', 'MB', 'GB', 'TB'];

/** A file size as people read it: 0 B, 12 KB, 3.4 MB. */
export function fileSize(bytes: number): string {
	let value = Math.max(0, bytes);
	let unit = 0;
	while (value >= 1024 && unit < UNITS.length - 1) {
		value /= 1024;
		unit++;
	}
	const shown = unit === 0 || value >= 10 ? Math.round(value) : Math.round(value * 10) / 10;
	return `${shown} ${UNITS[unit]}`;
}
