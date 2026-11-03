const AGENT_LOCALPART = /^@agent-/;

/** Reports whether a Matrix user id belongs to a Bulle agent account. */
export function isAgent(userId: string): boolean {
	return AGENT_LOCALPART.test(userId);
}

/** Reports whether a user id is the Prisme appservice bot, which users never address. */
export function isBot(userId: string): boolean {
	return userId.startsWith('@prisme:');
}
