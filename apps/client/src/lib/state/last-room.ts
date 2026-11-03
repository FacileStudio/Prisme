const KEY = 'prisme.last-room';

/** The room the user had open last, from this browser. */
export function lastRoom(): string | null {
	try {
		return localStorage.getItem(KEY);
	} catch {
		return null;
	}
}

/** Remembers the open room for the next visit. */
export function rememberRoom(roomId: string): void {
	try {
		localStorage.setItem(KEY, roomId);
	} catch {
		return;
	}
}
