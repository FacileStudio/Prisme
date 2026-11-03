import { describe, expect, it } from 'vitest';
import { composeMessage } from './compose';
import { escapeHtml, renderMarkdown, stripReplyFallback } from './markdown';
import { activeMention, filterMembers, insertMention, presentMentions } from './mentions';
import { isAgent } from './users';
import { continues, formatCount, roomBadge, sumBadges } from './unread';

describe('markdown', () => {
	it('returns null for plain text', () => {
		expect(renderMarkdown('hello there')).toBeNull();
		expect(renderMarkdown('two\nlines')).toBeNull();
	});

	it('renders inline formatting without a wrapping paragraph', () => {
		expect(renderMarkdown('**bold** and `code`')).toBe('<strong>bold</strong> and <code>code</code>');
	});

	it('renders fenced code blocks', () => {
		expect(renderMarkdown('```js\nlet a = 1;\n```')).toContain('<pre><code class="language-js">let a = 1;');
	});

	it('escapes raw html instead of passing it through', () => {
		const html = renderMarkdown('**x** <img src=x onerror=alert(1)>') ?? '';
		expect(html).not.toContain('<img');
		expect(html).toContain('&lt;img');
	});

	it('escapes special characters', () => {
		expect(escapeHtml(`<a href="x">'&'</a>`)).toBe('&lt;a href=&quot;x&quot;&gt;&#39;&amp;&#39;&lt;/a&gt;');
	});

	it('strips a reply fallback', () => {
		expect(stripReplyFallback('> <@a:b> hi\n> there\n\nanswer')).toBe('answer');
		expect(stripReplyFallback('no quote')).toBe('no quote');
	});
});

describe('mentions', () => {
	const members = [
		{ userId: '@alice:localhost', name: 'Alice', agent: false },
		{ userId: '@agent-alice-bulle:localhost', name: 'Bulle', agent: true },
		{ userId: '@bob:localhost', name: 'Bob', agent: false }
	];

	it('detects agents by localpart', () => {
		expect(isAgent('@agent-alice-bulle:localhost')).toBe(true);
		expect(isAgent('@alice:localhost')).toBe(false);
	});

	it('finds the query under the caret', () => {
		expect(activeMention('hi @bo', 6)).toEqual({ start: 3, query: 'bo' });
		expect(activeMention('@', 1)).toEqual({ start: 0, query: '' });
		expect(activeMention('mail a@b', 8)).toBeNull();
		expect(activeMention('hi @bob there', 13)).toBeNull();
	});

	it('ranks prefix matches before substring matches', () => {
		expect(filterMembers(members, 'b').map((m) => m.name)).toEqual(['Bob', 'Bulle']);
		expect(filterMembers(members, 'alice').map((m) => m.name)).toEqual(['Alice', 'Bulle']);
		expect(filterMembers(members, 'zz')).toEqual([]);
	});

	it('inserts the chosen name and moves the caret', () => {
		expect(insertMention('hi @bo', { start: 3, query: 'bo' }, 6, 'Bob')).toEqual({ text: 'hi @Bob ', caret: 8 });
	});

	it('drops mentions that were deleted from the text', () => {
		const chosen = [
			{ userId: '@bob:localhost', name: 'Bob' },
			{ userId: '@bob:localhost', name: 'Bob' },
			{ userId: '@alice:localhost', name: 'Alice' }
		];
		expect(presentMentions('hey @Bob', chosen)).toEqual([{ userId: '@bob:localhost', name: 'Bob' }]);
	});
});

describe('composeMessage', () => {
	it('sends plain text without html', () => {
		expect(composeMessage('  hello  ', [])).toEqual({ body: 'hello', html: null, userIds: [] });
	});

	it('turns mentions into matrix.to pills and m.mentions ids', () => {
		const out = composeMessage('hi @Bulle_bot do it', [{ userId: '@agent-x:localhost', name: 'Bulle_bot' }]);
		expect(out.body).toBe('hi @Bulle_bot do it');
		expect(out.html).toBe('hi <a href="https://matrix.to/#/@agent-x:localhost">@Bulle_bot</a> do it');
		expect(out.userIds).toEqual(['@agent-x:localhost']);
	});

	it('keeps markdown around mentions', () => {
		const out = composeMessage('**@Bob** look', [{ userId: '@bob:localhost', name: 'Bob' }]);
		expect(out.html).toBe('<strong><a href="https://matrix.to/#/@bob:localhost">@Bob</a></strong> look');
	});
});

describe('unread', () => {
	it('counts every DM message and only channel highlights', () => {
		expect(roomBadge({ total: 3, highlight: 0, muted: false, direct: true })).toEqual({ count: 3, strong: true, bold: true });
		expect(roomBadge({ total: 3, highlight: 1, muted: false, direct: false })).toEqual({ count: 1, strong: true, bold: true });
		expect(roomBadge({ total: 3, highlight: 0, muted: false, direct: false })).toEqual({ count: 0, strong: false, bold: true });
	});

	it('mutes the bold state but keeps mentions', () => {
		expect(roomBadge({ total: 5, highlight: 0, muted: true, direct: false }).bold).toBe(false);
		expect(roomBadge({ total: 5, highlight: 2, muted: true, direct: false }).count).toBe(2);
	});

	it('sums a space', () => {
		const rooms = [
			{ total: 2, highlight: 1, muted: false, direct: false },
			{ total: 4, highlight: 0, muted: true, direct: false }
		];
		expect(sumBadges(rooms)).toEqual({ count: 1, strong: true, bold: true });
		expect(sumBadges([])).toEqual({ count: 0, strong: false, bold: false });
	});

	it('caps counters and groups close messages', () => {
		expect(formatCount(120)).toBe('99+');
		expect(continues({ sender: 'a', ts: 0 }, { sender: 'a', ts: 60_000 })).toBe(true);
		expect(continues({ sender: 'a', ts: 0 }, { sender: 'b', ts: 1 })).toBe(false);
		expect(continues(undefined, { sender: 'a', ts: 1 })).toBe(false);
	});
});
