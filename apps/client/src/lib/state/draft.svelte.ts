import type { MatrixEvent } from 'matrix-js-sdk';

/** What a composer is doing besides writing a new message: replying or editing. */
export class DraftState {
	replyTo = $state.raw<MatrixEvent | null>(null);
	editing = $state.raw<MatrixEvent | null>(null);

	/** Starts a reply to an event. */
	reply(event: MatrixEvent) {
		this.editing = null;
		this.replyTo = event;
	}

	/** Starts editing one of the user's own messages. */
	edit(event: MatrixEvent) {
		this.replyTo = null;
		this.editing = event;
	}

	/** Drops the reply or edit. */
	clear() {
		this.replyTo = null;
		this.editing = null;
	}
}
