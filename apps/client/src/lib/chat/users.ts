import { t } from '$lib/i18n/index.svelte';

const AGENT_LOCALPART = /^@agent-/;

/**
 * Reports whether a Matrix user id belongs to an agent: listed in `known`
 * (from `GET /api/users/agents`), or a legacy `@agent-` account.
 */
export function isAgent(userId: string, known?: ReadonlyMap<string, unknown>): boolean {
	return AGENT_LOCALPART.test(userId) || !!known?.has(userId);
}

/** Reports whether a user id is the Prisme appservice bot, which users never address. */
export function isBot(userId: string): boolean {
	return userId.startsWith('@prisme:');
}

/** The `@name` people see for a Matrix user id, without the `:server` part. */
export function handle(userId: string): string {
	const colon = userId.indexOf(':');
	return colon === -1 ? userId : userId.slice(0, colon);
}

/** The sentence under the composer for the given names. */
export function typingSentence(names: string[]): string {
	if (names.length === 0) return '';
	if (names.length === 1) return t('typing.one', { a: names[0] });
	if (names.length === 2) return t('typing.two', { a: names[0], b: names[1] });
	return t('typing.many', { a: names[0], b: names[1], count: names.length - 2 });
}
