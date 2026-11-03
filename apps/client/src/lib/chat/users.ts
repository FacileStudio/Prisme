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

/** The sentence under the composer for the given names. */
export function typingSentence(names: string[]): string {
	if (names.length === 0) return '';
	if (names.length === 1) return `${names[0]} is typing…`;
	if (names.length === 2) return `${names[0]} and ${names[1]} are typing…`;
	return `${names[0]}, ${names[1]} and ${names.length - 2} more are typing…`;
}
