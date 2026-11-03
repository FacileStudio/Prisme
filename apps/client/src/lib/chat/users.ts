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
