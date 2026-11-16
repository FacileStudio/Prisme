import type { Component } from 'svelte';
import type { Room } from 'matrix-js-sdk';
import { glyphs } from '$lib/glyphs';
import { t } from '$lib/i18n/index.svelte';
import { app } from '$lib/state/app.svelte';
import type { ComposerState } from './composer.svelte';
import GifPanel from '$lib/components/room/GifPanel.svelte';

/** What the "+" menu hands an entry and the panel it opens. */
export type ComposerContext = {
	room: Room;
	threadId: string | null;
	composer: ComposerState;
	editing: boolean;
	pickFiles: () => void;
};

/** The props of a panel that opens above the composer. A panel calls `onclose` once it is done. */
export type ComposerPanelProps = { context: ComposerContext; onclose: () => void };

/**
 * One row of the composer's "+" menu. A row either runs an action (`run`) or opens a
 * panel above the composer (`panel`). `available` hides it, `active` shows it pressed.
 */
export type ComposerEntry = {
	id: string;
	icon: string;
	label: (context: ComposerContext) => string;
	available?: (context: ComposerContext) => boolean;
	active?: (context: ComposerContext) => boolean;
	run?: (context: ComposerContext) => void;
	panel?: Component<ComposerPanelProps>;
};

const idle = (context: ComposerContext) => !context.editing && !context.composer.sealed;

/**
 * The rows of the "+" menu, in the order they show. Add a feature by adding one entry
 * here; the composer renders the menu and the panel from this list alone.
 */
export const composerEntries: ComposerEntry[] = [
	{
		id: 'file',
		icon: glyphs.paperclip,
		label: () => t('composerMenu.file'),
		available: idle,
		run: (context) => context.pickFiles()
	},
	{
		id: 'capsule',
		icon: glyphs.capsule,
		label: (context) => (context.composer.sealed ? t('composerMenu.capsuleOff') : t('composerMenu.capsule')),
		available: (context) => !context.editing && context.composer.files.length === 0,
		active: (context) => context.composer.sealed,
		run: (context) => (context.composer.sealed = !context.composer.sealed)
	},
	{
		id: 'gif',
		icon: glyphs.gif,
		label: () => t('composerMenu.gif'),
		available: (context) => idle(context) && !!app.session?.gifs_enabled,
		panel: GifPanel
	}
];

/** The rows the menu shows for this composer right now. */
export function visibleEntries(context: ComposerContext, entries: ComposerEntry[] = composerEntries): ComposerEntry[] {
	return entries.filter((entry) => entry.available?.(context) ?? true);
}
