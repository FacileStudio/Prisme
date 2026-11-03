const auth = { base: '', token: '' };
const cache = new Map<string, Promise<string>>();
const LIMIT = 150;

/** Stores the homeserver and token the authenticated media routes need. */
export function setMediaToken(base: string, token: string): void {
	auth.base = base.replace(/\/$/, '');
	auth.token = token;
}

/** Builds the authenticated media URL for an `mxc://` URI, or a thumbnail of it. */
export function mediaUrl(mxc: string, size?: number): string | null {
	const match = /^mxc:\/\/([^/]+)\/([^/?#]+)$/.exec(mxc);
	if (!match) return null;
	const [, server, id] = match;
	if (!size) return `${auth.base}/_matrix/client/v1/media/download/${server}/${id}`;
	return `${auth.base}/_matrix/client/v1/media/thumbnail/${server}/${id}?width=${size}&height=${size}&method=scale`;
}

function evict() {
	while (cache.size > LIMIT) {
		const [key, url] = cache.entries().next().value as [string, Promise<string>];
		cache.delete(key);
		void url.then((blob) => URL.revokeObjectURL(blob)).catch(() => undefined);
	}
}

async function fetchBlob(url: string): Promise<string> {
	const response = await fetch(url, { headers: { Authorization: `Bearer ${auth.token}` } });
	if (!response.ok) throw new Error(`media ${response.status}`);
	return URL.createObjectURL(await response.blob());
}

/**
 * Resolves an `mxc://` URI to a blob URL, fetched with the bearer token because
 * Synapse serves media only to authenticated requests and `<img>` cannot send one.
 */
export function loadMedia(mxc: string, size?: number): Promise<string> {
	const url = mediaUrl(mxc, size);
	if (!url) return Promise.reject(new Error('not an mxc uri'));
	const hit = cache.get(url);
	if (hit) {
		cache.delete(url);
		cache.set(url, hit);
		return hit;
	}
	const pending = fetchBlob(url);
	pending.catch(() => cache.delete(url));
	cache.set(url, pending);
	evict();
	return pending;
}
