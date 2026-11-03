export type Panel = 'thread' | 'members' | 'search' | null;
export type Dialog =
	| 'create-space'
	| 'create-channel'
	| 'invite'
	| 'new-dm'
	| 'add-agent'
	| 'rename-room'
	| 'rename-space'
	| 'settings'
	| null;

export const HOME = 'home';

class UiState {
	spaceId = $state<string>(HOME);
	panel = $state<Panel>(null);
	threadId = $state<string | null>(null);
	drawer = $state(false);
	dialog = $state<Dialog>(null);

	/** Opens a thread in the side panel. */
	openThread(rootId: string) {
		this.threadId = rootId;
		this.panel = 'thread';
	}

	/** Toggles a side panel, closing it when it is already open. */
	toggle(panel: Exclude<Panel, null>) {
		this.panel = this.panel === panel ? null : panel;
	}

	/** Closes the side panel. */
	closePanel() {
		this.panel = null;
		this.threadId = null;
	}
}

export const ui = new UiState();
