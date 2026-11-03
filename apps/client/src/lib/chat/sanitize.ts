import DOMPurify from 'dompurify';
import { handle } from './users';

const ALLOWED_TAGS = [
	'a', 'b', 'blockquote', 'br', 'code', 'del', 'details', 'div', 'em', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
	'hr', 'i', 'li', 'ol', 'p', 'pre', 's', 'span', 'strong', 'sub', 'summary', 'sup', 'table', 'tbody',
	'td', 'th', 'thead', 'tr', 'u', 'ul'
];

const PILL = /^https:\/\/matrix\.to\/#\/(@[^/?]+)/;

function decorateLink(node: Element, me: string) {
	const href = node.getAttribute('href') ?? '';
	const pill = PILL.exec(href);
	if (pill) {
		const userId = decodeURIComponent(pill[1]);
		node.setAttribute('class', userId === me ? 'mention-pill mention-me' : 'mention-pill');
		node.setAttribute('data-user', userId);
		node.textContent = handle(userId);
		return;
	}
	node.setAttribute('target', '_blank');
	node.setAttribute('rel', 'noopener noreferrer');
}

/**
 * Sanitizes a received `formatted_body` to the Matrix-safe tag set, drops the
 * `mx-reply` fallback, and shows matrix.to user links as `@username` mention pills.
 */
export function sanitizeMessage(html: string, me: string): string {
	const withoutReply = html.replace(/<mx-reply>[\s\S]*?<\/mx-reply>/gi, '');
	DOMPurify.addHook('afterSanitizeAttributes', (node) => {
		if (node.tagName === 'A') decorateLink(node, me);
	});
	try {
		return DOMPurify.sanitize(withoutReply, {
			ALLOWED_TAGS,
			ALLOWED_ATTR: ['href', 'class', 'data-mx-color', 'start', 'colspan', 'rowspan'],
			ALLOWED_URI_REGEXP: /^(?:https?|mailto):/i
		});
	} finally {
		DOMPurify.removeHook('afterSanitizeAttributes');
	}
}
