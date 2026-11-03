export type Panel = 'thread' | 'members' | 'search' | null;
export type Dialog =
	| 'create-space'
	| 'create-channel'
	| 'invite'
	| 'add-agent'
	| 'rename-room'
	| 'rename-space'
	| null;

export const HOME = 'home';
export const MANAGE_SPACES_HREF = 'https://sso.facile.studio/spaces';
export const MANAGE_SPACES_LABEL = 'Manage shared spaces';

class UiState {
	spaceId = $state<string>(HOME);
	panel = $state<Panel>(null);
	threadId = $state<string | null>(null);
	drawer = $state(false);
	railCollapsed = $state(false);
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
