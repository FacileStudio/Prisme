import { tick } from 'svelte';
import type { ComposerContext, ComposerEntry } from './composer-menu';

type Targets = {
	plus: () => HTMLElement | undefined;
	smiley: () => HTMLElement | undefined;
	input: () => HTMLElement | undefined;
};

/**
 * The composer's popovers: the "+" menu, the panel one of its entries opened, and the
 * emoji picker. At most one is open, and closing one hands focus back where it belongs.
 */
export class ComposerPopups {
	menu = $state(false);
	emoji = $state(false);
	panel = $state<ComposerEntry | null>(null);

	constructor(private targets: Targets) {}

	/** Opens or closes the "+" menu. */
	toggleMenu() {
		const next = !this.menu;
		this.closeAll();
		this.menu = next;
	}

	/** Opens or closes the emoji picker. */
	toggleEmoji() {
		const next = !this.emoji;
		this.closeAll();
		this.emoji = next;
	}

	/** Closes the "+" menu, back to its button when the keyboard closed it. */
	closeMenu(refocus: boolean) {
		this.menu = false;
		if (refocus) this.targets.plus()?.focus();
	}

	/** Closes the emoji picker, back to its button when the keyboard closed it. */
	closeEmoji(refocus: boolean) {
		this.emoji = false;
		if (refocus) this.targets.smiley()?.focus();
	}

	/** Runs a menu entry, or opens its panel. */
	choose(entry: ComposerEntry, context: ComposerContext) {
		this.menu = false;
		if (entry.panel) this.panel = entry;
		else {
			entry.run?.(context);
			this.refocus();
		}
	}

	/** Closes the open panel and returns to the text field. */
	closePanel() {
		this.panel = null;
		this.refocus();
	}

	private refocus() {
		void tick().then(() => this.targets.input()?.focus());
	}

	private closeAll() {
		this.menu = false;
		this.emoji = false;
		this.panel = null;
	}
}
