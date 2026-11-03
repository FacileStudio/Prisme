import { Marked } from 'marked';

const ESCAPES: Record<string, string> = { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' };

/** Escapes text for safe inclusion in HTML. */
export function escapeHtml(text: string): string {
	return text.replace(/[&<>"']/g, (char) => ESCAPES[char]);
}

const marked = new Marked({
	gfm: true,
	breaks: true,
	renderer: {
		html: ({ text }) => escapeHtml(text)
	}
});

/**
 * Renders composer markdown to Matrix `formatted_body` HTML. Returns null when
 * the text has no formatting, so a plain message is sent without HTML.
 */
export function renderMarkdown(text: string): string | null {
	const html = (marked.parse(text, { async: false }) as string).trim();
	const plain = `<p>${escapeHtml(text).replace(/\n/g, '<br>')}</p>`;
	if (html === plain) return null;
	return html.replace(/^<p>([\s\S]*)<\/p>$/, (whole, inner: string) => (inner.includes('<p>') ? whole : inner));
}

/** Strips the reply fallback quote that older clients prepend to a reply body. */
export function stripReplyFallback(body: string): string {
	const lines = body.split('\n');
	let index = 0;
	while (index < lines.length && lines[index].startsWith('> ')) index++;
	if (index === 0) return body;
	return lines.slice(index).join('\n').replace(/^\n/, '');
}
